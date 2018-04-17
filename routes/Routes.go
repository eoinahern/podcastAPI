package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/eoinahern/podcastAPI/models"
	"github.com/eoinahern/podcastAPI/repository"
	"github.com/eoinahern/podcastAPI/util"
	"github.com/eoinahern/podcastAPI/validation"

	"github.com/gorilla/mux"
)

//RegisterHandler : register user initially
type RegisterHandler struct {
	EmailValidator  *validation.EmailValidation
	DB              repository.UserDBInt
	MailHelper      util.MailRequestInt
	PassEncryptUtil *util.PasswordEncryptUtil
}

//ConfirmRegistrationHandler : confirm registration via email get request
type ConfirmRegistrationHandler struct {
	DB repository.UserDBInt
}

//CreateSessionHandler : create a session and return jwt token
type CreateSessionHandler struct {
	DB              repository.UserDBInt
	PassEncryptUtil util.PasswordEncryptUtilInt
	JwtTokenUtil    *util.JwtTokenUtil
}

//CreatePodcastHandler : allows user to create a podcast
type CreatePodcastHandler struct {
	PodcastDB    repository.PodcastDBInt
	FileHelper   util.FileHelperUtilInt
	BaseLocation string
}

//GetPodcastsHandler : get all podcasts
type GetPodcastsHandler struct {
	UserDB    repository.UserDBInt
	PodcastDB repository.PodcastDBInt
}

//GetEpisodesHandler : all episodes associated with specific podcast
type GetEpisodesHandler struct {
	UserDB    repository.UserDBInt
	EpisodeDB repository.EpisodeDBInt
}

//DownloadEpisodeHandler : download a specific episode data
type DownloadEpisodeHandler struct {
	EpisodeDB    *repository.EpisodeDB
	BaseLocation string
}

//UploadEpisodeHandler : allows admin of a podcast to upload an episode file
type UploadEpisodeHandler struct {
	//credentials. then upload to network
	UserDB    repository.UserDBInt
	EpisodeDB repository.EpisodeDBInt
	PodcastDB repository.PodcastDBInt
}

//DeleteEpisodeHandler delete episode from specific podcast. Admin use
type DeleteEpisodeHandler struct {
	UserDB    repository.UserDBInt
	PodcastDB *repository.PodcastDB
}

//vars
var tokenErr = []byte(`{ "error" : "problem with token"}`)
var internalErr = []byte(`{ "error" : "internal error"}`)

const notAllowedErrStr string = "method not allowed"

func (r *RegisterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	//needs updating.
	//1 . create use once token.
	//2. add user to DB but unverified.
	//3. send confirmation emailValidator
	//4. new route to handle reg verification

	decoder := json.NewDecoder(req.Body)
	var user models.User
	err := decoder.Decode(&user)

	if err != nil {
		panic(err)
	}

	if len(user.UserName) == 0 || len(user.Password) == 0 {
		http.Error(w, "incorrect params", http.StatusBadRequest)
		return
	}

	if isValidEmail := r.EmailValidator.CheckEmailValid(user.UserName); !isValidEmail {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	if r.DB.CheckExist(user.UserName) {
		http.Error(w, http.StatusText(31), http.StatusConflict)
		return
	}

	user.Password = r.PassEncryptUtil.Encrypt(user.Password)
	user.RegToken = util.GenerateRandomToken()

	r.DB.Insert(&user)

	//send automated email with email link#

	r.MailHelper.SetBodyParams(&models.TemplateParams{User: user.UserName, Token: user.RegToken})
	r.MailHelper.SetToID(user.UserName)
	_, err = r.MailHelper.SendMail()

	if err != nil {
		log.Println("error sending automated mail")
		log.Println(err)
		http.Error(w, http.StatusText(51), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	msg := &models.Message{Message: fmt.Sprintf("registration confirmation email sent to %s", user.UserName)}
	resp, _ := json.Marshal(msg)
	w.Write(resp)
}

/**
* TODO: this needs to change from get request. need to pass required params in body
* could use put, patch? OR im opting for POST with body and a url param operation=confirm
* to me this seems like its legal a GET should be resource locator and shouldnt be
* passing data in GET url according to API rules.
**/

func (c *ConfirmRegistrationHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()
	user := params.Get("user")
	token := params.Get("token")

	w.Header().Set("Content-Type", "text/html")

	if c.DB.ValidateUserPlusRegToken(user, token) {
		c.DB.SetVerified(user, token)
		w.Write([]byte(fmt.Sprintf("<h1>  user %s registration confirmed<h1>", user)))
		return
	}

	w.Write([]byte("<h1> problem verifying user? <h1>"))
}

func (c *CreateSessionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var user models.User
	err := decoder.Decode(&user)

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	dbUser := c.DB.GetUser(user.UserName)
	w.Header().Set("Content-Type", "application/json")

	if c.PassEncryptUtil.CheckSame(dbUser.Password, user.Password) {
		user.Token = c.JwtTokenUtil.CreateToken(user.UserName)
		jsonUser, _ := json.Marshal(user)
		w.Write(jsonUser)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error" : "incorrect pass"}`))
	}
}

//create a podcast entry and folder on server.

func (c *CreatePodcastHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	//1. authorize user .....
	//2. create user folder if not exists. create podcast folder.
	//3. store data in DB about podcast
	//3. return success

	w.Header().Set("Content-Type", "application/json")

	var podcast models.Podcast
	err := json.NewDecoder(req.Body).Decode(&podcast)

	if err != nil {
		http.Error(w, http.StatusText(51), http.StatusInternalServerError)
		return
	}

	podcastname := req.URL.Query().Get("podcastname")

	if len(podcastname) == 0 {
		http.Error(w, http.StatusText(22), http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("%s/%d/%s", c.BaseLocation, podcast.PodcastID, podcastname)

	if !c.FileHelper.CheckDirFileExists(path) {
		c.FileHelper.CreateDir(path)
		podcast.Location = path
		podcast.Name = podcastname
		err = c.PodcastDB.CreatePodcast(&podcast)

		if err != nil {
			http.Error(w, http.StatusText(51), http.StatusInternalServerError)
			return
		}

		mpod, _ := json.Marshal(podcast)
		w.Write(mpod)
	}
}

// get a list of the most popular podcasts and return to the users
// in json format

func (g *GetPodcastsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	podcasts := g.PodcastDB.GetAll()
	podcastsMarshaled, err := json.Marshal(podcasts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(internalErr)
	} else {
		w.Write(podcastsMarshaled)
	}
}

/**
*	get episodes from a specific podcast
* by podcast id!!!
**/

func (g *GetEpisodesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	podcastid, err := strconv.Atoi(req.URL.Query().Get("podcastid"))

	if err != nil {
		http.Error(w, http.StatusText(22), http.StatusBadRequest)
		return
	}

	episodes := g.EpisodeDB.GetAllEpisodes(podcastid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&episodes)
}

func (g *DownloadEpisodeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	requestParams := mux.Vars(req)

	podcastID := requestParams["podcastid"]
	podcastName := requestParams["podcastname"]
	podcastFileName := requestParams["podcastfilename"]

	if len(podcastID) == 0 || len(podcastName) == 0 || len(podcastFileName) == 0 {
		http.Error(w, "unrecognised", http.StatusBadRequest)
		return
	}

	podlocation := fmt.Sprintf("%s/%s/%s/%s", g.BaseLocation, podcastID, podcastName, podcastFileName)
	filedata, err := ioutil.ReadFile(podlocation)

	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Write(filedata)

}

func (e *UploadEpisodeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var episode models.Episode
	file, fh, fileErr := req.FormFile("namefile")
	sepisode := req.FormValue("data")
	podcastname := req.URL.Query().Get("podcast")
	err := json.Unmarshal([]byte(sepisode), &episode)

	if len(sepisode) == 0 || err != nil || fileErr != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	podcast := e.PodcastDB.CheckPodcastCreated(episode.PodID, podcastname)

	if len(podcast.Name) == 0 {
		http.Error(w, "unknown podcast", http.StatusInternalServerError)
		return
	}

	splitname := strings.Split(fh.Filename, ".")
	ext := splitname[len(splitname)-1]

	if strings.Compare(ext, "mp3") != 0 {
		http.Error(w, "wrong file type", http.StatusInternalServerError)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	lastepisode := e.EpisodeDB.GetLastEpisode()
	filelocation := fmt.Sprintf("%s/%d.%s", podcast.Location, lastepisode.EpisodeID+1, "mp3")
	episode.URL = filelocation

	e.EpisodeDB.AddEpisode(episode)
	ioutil.WriteFile(filelocation, fileBytes, os.ModePerm)
	e.PodcastDB.UpdatePodcastNumEpisodes(podcast.PodcastID)
}
