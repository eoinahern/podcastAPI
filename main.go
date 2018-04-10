package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

func main() {

	file, err := os.Open("config.json")

	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	config := map[string]models.Config{}
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

	/*passEncryptUtil := &util.PasswordEncryptUtil{}
	emailValidator := &validation.EmailValidation{}
	fileHelperUtil := &util.FileHelperUtil{}
	userDB := &repository.UserDB{DB: prodDB}
	episodeDB := &repository.EpisodeDB{DB: prodDB}
	podcastDB := &repository.PodcastDB{DB: prodDB}
	jwtTokenUtil := &util.JwtTokenUtil{SigningKey: prodConf.SigningKey, DB: userDB}
	regMailHelper := &util.MailRequest{SenderId: "mypodcastapi@gmail.com", BodyLocation: "view/templates/regMailTemplate.html"}*/

	//db.AutoMigrate(&models.User{}, &models.Podcast{}, &models.Episode{})
	//db.Model(&models.Podcast{}).AddForeignKey("user_email", "users(user_name)", "CASCADE", "CASCADE")
	//db.Model(&models.Episode{}).AddForeignKey("pod_id", "podcasts(podcast_id)", "CASCADE", "CASCADE")

	defer prodDB.Close()
	defer prodDB.Close()

	router := mux.NewRouter()

	setUpProduction(router, prodDB, prodConf.SigningKey)
	setUpDebug(router, debugDB)

	/*router.Handle("/register", &routes.RegisterHandler{EmailValidator: emailValidator, MailHelper: regMailHelper, DB: userDB, PassEncryptUtil: passEncryptUtil}).Methods(http.MethodPost)
	router.Handle("/confirm", &routes.ConfirmRegistrationHandler{DB: userDB}).Methods(http.MethodPost)
	router.Handle("/session", &routes.CreateSessionHandler{DB: userDB, JwtTokenUtil: jwtTokenUtil, PassEncryptUtil: passEncryptUtil}).Methods(http.MethodPost)
	router.Handle("/podcasts", middleware.Adapt(&routes.GetPodcastsHandler{UserDB: userDB, PodcastDB: podcastDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/episodes", middleware.Adapt(&routes.GetEpisodesHandler{UserDB: userDB, EpisodeDB: episodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/episodes/{podcastid}/{podcastname}/{podcastfilename}", middleware.Adapt(&routes.DownloadEpisodeHandler{EpisodeDB: episodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodGet)
	router.Handle("/podcasts", middleware.Adapt(&routes.CreatePodcastHandler{PodcastDB: podcastDB, FileHelper: fileHelperUtil}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodPost)
	router.Handle("/episodes", middleware.Adapt(&routes.UploadEpisodeHandler{UserDB: userDB, PodcastDB: podcastDB, EpisodeDB: episodeDB}, middleware.AuthMiddlewareInit(jwtTokenUtil))).Methods(http.MethodPost)
	*/

	http.ListenAndServe(":8080", router)
	//http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
}

func setUpProduction(router *mux.Router, prodDB *sql.DB, signingKey string) {

	passEncryptUtil := &util.PasswordEncryptUtil{}
	emailValidator := &validation.EmailValidation{}
	fileHelperUtil := &util.FileHelperUtil{}
	userDB := &repository.UserDB{DB: prodDB}
	episodeDB := &repository.EpisodeDB{DB: prodDB}
	podcastDB := &repository.PodcastDB{DB: prodDB}
	jwtTokenUtil := &util.JwtTokenUtil{SigningKey: signingKey, DB: userDB}
	regMailHelper := &util.MailRequest{SenderId: "mypodcastapi@gmail.com", BodyLocation: "view/templates/regMailTemplate.html"}

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

}

func seed() {

}
