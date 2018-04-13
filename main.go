package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/eoinahern/podcastAPI/validation"

	"github.com/eoinahern/podcastAPI/middleware"
	"github.com/eoinahern/podcastAPI/models"
	"github.com/eoinahern/podcastAPI/repository"
	"github.com/eoinahern/podcastAPI/routes"
	"github.com/eoinahern/podcastAPI/util"

	"github.com/gorilla/mux"
)

var emailValidator *validation.EmailValidation
var passEncryptUtil *util.PasswordEncryptUtil
var fileHelperUtil *util.FileHelperUtil
var regMailHelper *util.MailRequest
var debugUserDB *repository.UserDB
var debugPodcastDB *repository.PodcastDB
var debugEpisodeDB *repository.EpisodeDB
var wg sync.WaitGroup

const seedFileLocation string = "config/seedData.json"

func main() {

	wg.Add(1)
	file, err := os.Open("config/config.json")

	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	var config map[string]models.Config
	decoder.Decode(&config)

	prodConf := config["production"]
	debugConf := config["debug"]

	prodConfDataSource := fmt.Sprintf("%s:%s@/%s", prodConf.User, prodConf.Password, prodConf.Schema)
	debugConfDataSource := fmt.Sprintf("%s:%s@/%s", debugConf.User, debugConf.Password, debugConf.Schema)

	prodDB, err := sql.Open("mysql", prodConfDataSource)
	debugDB, err := sql.Open("mysql", debugConfDataSource)

	if pingErr := prodDB.Ping(); pingErr != nil {
		panic("error connecting db " + err.Error())
	}

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	//create dependencies
	passEncryptUtil = &util.PasswordEncryptUtil{}
	emailValidator = &validation.EmailValidation{}
	fileHelperUtil = &util.FileHelperUtil{}
	regMailHelper = &util.MailRequest{SenderId: "mypodcastapi@gmail.com", BodyLocation: "view/templates/regMailTemplate.html"}

	defer prodDB.Close()
	defer debugDB.Close()

	router := mux.NewRouter()

	setUpDebug(router, debugDB)
	wg.Wait()
	setUpProduction(router, prodDB, prodConf.SigningKey)
	http.ListenAndServe(":8080", router)
	//http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
}

func setUpProduction(router *mux.Router, prodDB *sql.DB, signingKey string) {

	var userDB repository.UserDBInt = &repository.UserDB{DB: prodDB}
	episodeDB := &repository.EpisodeDB{DB: prodDB}
	podcastDB := &repository.PodcastDB{DB: prodDB}
	jwtTokenUtil := &util.JwtTokenUtil{SigningKey: signingKey, DB: userDB}

	router.Handle("/register", &routes.RegisterHandler{EmailValidator: emailValidator, MailHelper: regMailHelper, DB: userDB, PassEncryptUtil: passEncryptUtil}).Methods(http.MethodPost)
	router.Handle("/confirm", &routes.ConfirmRegistrationHandler{DB: userDB}).Methods(http.MethodPost)
	router.Handle("/session", &routes.CreateSessionHandler{DB: userDB, JwtTokenUtil: jwtTokenUtil, PassEncryptUtil: passEncryptUtil}).Methods(http.MethodPost)
	router.Handle("/podcasts", middleware.Adapt(&routes.GetPodcastsHandler{UserDB: userDB, PodcastDB: podcastDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/episodes", middleware.Adapt(&routes.GetEpisodesHandler{UserDB: userDB, EpisodeDB: episodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/episodes/{podcastid}/{podcastname}/{podcastfilename}", middleware.Adapt(&routes.DownloadEpisodeHandler{EpisodeDB: episodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/podcasts", middleware.Adapt(&routes.CreatePodcastHandler{PodcastDB: podcastDB, FileHelper: fileHelperUtil}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodPost)
	router.Handle("/episodes", middleware.Adapt(&routes.UploadEpisodeHandler{UserDB: userDB, PodcastDB: podcastDB, EpisodeDB: episodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodPost)

}

func setUpDebug(router *mux.Router, debugDB *sql.DB) {

	debugUserDB = &repository.UserDB{DB: debugDB}
	debugEpisodeDB = &repository.EpisodeDB{DB: debugDB}
	debugPodcastDB = &repository.PodcastDB{DB: debugDB}
	jwtTokenUtil := &util.JwtTokenUtil{SigningKey: "1234", DB: debugUserDB}

	router.Handle("/debug/register", &routes.RegisterHandler{EmailValidator: emailValidator, MailHelper: regMailHelper, DB: debugUserDB, PassEncryptUtil: passEncryptUtil}).Methods(http.MethodPost)
	router.Handle("/debug/confirm", &routes.ConfirmRegistrationHandler{DB: debugUserDB}).Methods(http.MethodPost)
	router.Handle("/debug/session", &routes.CreateSessionHandler{DB: debugUserDB, JwtTokenUtil: jwtTokenUtil, PassEncryptUtil: passEncryptUtil}).Methods(http.MethodPost)
	router.Handle("/debug/podcasts", middleware.Adapt(&routes.GetPodcastsHandler{UserDB: debugUserDB, PodcastDB: debugPodcastDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/debug/episodes", middleware.Adapt(&routes.GetEpisodesHandler{UserDB: debugUserDB, EpisodeDB: debugEpisodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/debug/episodes/{podcastid}/{podcastname}/{podcastfilename}", middleware.Adapt(&routes.DownloadEpisodeHandler{EpisodeDB: debugEpisodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/debug/podcasts", middleware.Adapt(&routes.CreatePodcastHandler{PodcastDB: debugPodcastDB, FileHelper: fileHelperUtil}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodPost)
	router.Handle("/debug/episodes", middleware.Adapt(&routes.UploadEpisodeHandler{UserDB: debugUserDB, PodcastDB: debugPodcastDB, EpisodeDB: debugEpisodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodPost)

	go seed()
}

func seed() {

	fmt.Println("seeding Database!")
	defer wg.Done()

	var seedData models.SeedData
	userRows := debugUserDB.CountRows()
	podcastRows := debugPodcastDB.CountRows()
	episodeRows := debugEpisodeDB.CountRows()

	if userRows == 0 || podcastRows == 0 || episodeRows == 0 {
		file, _ := os.Open(seedFileLocation)
		json.NewDecoder(file).Decode(&seedData)
		defer file.Close()
	}

	if userRows == 0 {
		debugUserDB.Insert(&seedData.User)
	}

	if podcastRows == 0 {
		for _, podcast := range seedData.Podcasts {
			debugPodcastDB.CreatePodcast(podcast)
		}
	}

	if episodeRows == 0 {
		for _, episode := range seedData.Episodes {
			debugEpisodeDB.AddEpisode(*episode)
		}
	}

	fmt.Println("seed Debug db complete!")
}
