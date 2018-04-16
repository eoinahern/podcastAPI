package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"

	"github.com/eoinahern/podcastAPI/mocks"
	"github.com/eoinahern/podcastAPI/models"
	"github.com/eoinahern/podcastAPI/util"
	"github.com/eoinahern/podcastAPI/validation"
)

const host string = "localhost"

func TestRegisterHandler(t *testing.T) {

	registerHandler := &RegisterHandler{
		EmailValidator:  &validation.EmailValidation{},
		DB:              &mocks.MockUserDB{},
		MailHelper:      &mocks.MockMailRequest{},
		PassEncryptUtil: &util.PasswordEncryptUtil{},
	}

	responseWriter := httptest.NewRecorder()
	user, _ := json.Marshal(models.User{UserName: "eoin@yahoo.co.uk", Password: "hellothere"})
	request, err := http.NewRequest("POST", host, bytes.NewReader(user))

	if err != nil {
		t.Error(err)
	}

	registerHandler.ServeHTTP(responseWriter, request)
	assert.Equal(t, 200, responseWriter.Code)
	assert.Equal(t, "application/json", responseWriter.Header().Get("Content-Type"))

	// call with empty params
	user, _ = json.Marshal(models.User{UserName: "", Password: ""})
	ReqNoPass, _ := http.NewRequest(http.MethodPost, host, bytes.NewReader(user))
	responseWriter = httptest.NewRecorder()

	registerHandler.ServeHTTP(responseWriter, ReqNoPass)
	assert.Equal(t, "incorrect params\n", responseWriter.Body.String())
	assert.Equal(t, http.StatusBadRequest, responseWriter.Code)

}

func TestConfirmRegistration(t *testing.T) {

	regHandler := &ConfirmRegistrationHandler{DB: &mocks.MockUserDB{}}

	request, err := http.NewRequest(http.MethodPost, host, nil)
	responseWriter := httptest.NewRecorder()

	if err != nil {
		t.Error(err)
	}

	regHandler.ServeHTTP(responseWriter, request)
	assert.Equal(t, http.StatusOK, responseWriter.Code)
	assert.Equal(t, "text/html", responseWriter.Header().Get("Content-Type"))
	assert.Equal(t, "<h1> problem verifying user? <h1>", responseWriter.Body.String())

	// return reg confirm pass
	requestPass, _ := http.NewRequest(http.MethodPost, "localhost?user=eoin@yahoo.com&token=1234", nil)
	responseWriterPass := httptest.NewRecorder()

	regHandler.ServeHTTP(responseWriterPass, requestPass)
	assert.Equal(t, "<h1>  user eoin@yahoo.com registration confirmed<h1>", responseWriterPass.Body.String())

}

func TestCreateSession(t *testing.T) {

	createSessionHandler := &CreateSessionHandler{DB: &mocks.MockUserDB{}, PassEncryptUtil: &mocks.MockPasswordEncryptUitl{}, JwtTokenUtil: &util.JwtTokenUtil{}}

	user, _ := json.Marshal(models.User{UserName: "eoin@yahoo.co.uk", Password: "hellothere"})
	request, err := http.NewRequest(http.MethodPost, host, bytes.NewReader(user))

	if err != nil {
		t.Error(err)
	}
	respWriter := httptest.NewRecorder()
	createSessionHandler.ServeHTTP(respWriter, request)

	var returnedUser models.User
	json.NewDecoder(respWriter.Body).Decode(&returnedUser)

	assert.Equal(t, "eoin@yahoo.co.uk", returnedUser.UserName)

}

func TestGetPodcasts(t *testing.T) {

	getPodcastsHandler := &GetPodcastsHandler{UserDB: &mocks.MockUserDB{}, PodcastDB: &mocks.MockPodcastDB{}}

	request, err := http.NewRequest(http.MethodGet, host, nil)
	respWriter := httptest.NewRecorder()

	if err != nil {
		t.Error(err)
	}

	getPodcastsHandler.ServeHTTP(respWriter, request)

	var podcasts []models.SecurePodcast
	json.NewDecoder(respWriter.Body).Decode(&podcasts)

	assert.Equal(t, 2, len(podcasts))
	assert.Equal(t, 2, podcasts[0].EpisodeNum)
	assert.Equal(t, 2, podcasts[1].EpisodeNum)

}

func TestCreatePodcast(t *testing.T) {

}

func TestGetEPisode(t *testing.T) {

	getEpisodeHandler := &GetEpisodesHandler{}

	respWriter := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, host, nil)

	if err != nil {
		t.Error(err)
	}

	getEpisodeHandler.ServeHTTP(respWriter, request)

}

func TestDownloadEpisode(t *testing.T) {

	downloadEpisodeHandler := &DownloadEpisodeHandler{BaseLocation: "../debug_files"}
	reqURL := "http://localhost:8080/episodes/{podcastid}/{podcastname}/{podcastfilename}"
	request, err := http.NewRequest(http.MethodGet, reqURL, nil)

	if err != nil {
		t.Error(err)
	}
	respWriter := httptest.NewRecorder()

	downloadEpisodeHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, respWriter.Code, http.StatusBadRequest)

	request = mux.SetURLVars(request, map[string]string{"podcastid": "pod", "podcastname": "name", "podcastfilename": "filename"})
	respWriter = httptest.NewRecorder()

	downloadEpisodeHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusNotFound, respWriter.Code)

	request = mux.SetURLVars(request, map[string]string{"podcastid": "test", "podcastname": "mypod", "podcastfilename": "sample.mp3"})
	respWriter = httptest.NewRecorder()

	downloadEpisodeHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusOK, respWriter.Code)
	assert.Equal(t, "audio/mpeg", respWriter.Header().Get("Content-Type"))

}

func TestUploadEpisode(t *testing.T) {

}
