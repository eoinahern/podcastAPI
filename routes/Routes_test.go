package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

	request, err := http.NewRequest(http.MethodGet, "http://localhost/debug/podcasts?limit=20&offset=0", nil)
	respWriter := httptest.NewRecorder()

	if err != nil {
		t.Error(err)
	}

	getPodcastsHandler.ServeHTTP(respWriter, request)

	var podcastsPage *models.PodcastPage
	json.NewDecoder(respWriter.Body).Decode(&podcastsPage)

	assert.Equal(t, 2, len(podcastsPage.Data))
	assert.Equal(t, 2, podcastsPage.Data[0].EpisodeNum)
	assert.Equal(t, 2, podcastsPage.Data[1].EpisodeNum)

}

func TestCreatePodcast(t *testing.T) {

	createPodcastHandler := &CreatePodcastHandler{PodcastDB: &mocks.MockPodcastDB{}, FileHelper: &mocks.MockFileHelperUtil{}, BaseLocation: "../debug_files"}

	respWriter := httptest.NewRecorder()
	podcast := &models.Podcast{PodcastID: 1, UserEmail: "eoin@yahoo.co.uk", Name: "podcast", Icon: "", Details: "a podcast"}
	encodedPodcast, _ := json.Marshal(podcast)
	request, err := http.NewRequest(http.MethodPost, "localhost/podcasts?podcastname=podcast", bytes.NewBuffer(encodedPodcast))

	if err != nil {
		t.Error(err)
	}

	createPodcastHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusOK, respWriter.Code)

	request, _ = http.NewRequest(http.MethodPost, "localhost/podcasts?podcastname=pod", bytes.NewBuffer(encodedPodcast))
	respWriter = httptest.NewRecorder()

	createPodcastHandler.ServeHTTP(respWriter, request)
	var podcastReturned models.Podcast
	json.NewDecoder(respWriter.Body).Decode(&podcastReturned)

	assert.Equal(t, http.StatusOK, respWriter.Code)
	assert.Equal(t, "pod", podcastReturned.Name)

	request, _ = http.NewRequest(http.MethodPost, "localhost/podcasts?podcastname=", bytes.NewBuffer([]byte("{}")))
	respWriter = httptest.NewRecorder()

	createPodcastHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

	// test with no request body
	request, _ = http.NewRequest(http.MethodPost, "localhost/podcasts?podcastname=podcast", bytes.NewBuffer([]byte("")))
	respWriter = httptest.NewRecorder()

	createPodcastHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusInternalServerError, respWriter.Code)

}

func TestGetEpisode(t *testing.T) {

	getEpisodeHandler := &GetEpisodesHandler{UserDB: &mocks.MockUserDB{}, EpisodeDB: &mocks.MockEpisodeDB{}}

	respWriter := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "http://localhost/debug/episodes?podcastid=1&limit=10&offset=60", nil)

	if err != nil {
		t.Error(err)
	}

	getEpisodeHandler.ServeHTTP(respWriter, request)

	var page models.EpisodePage
	json.NewDecoder(respWriter.Body).Decode(&page)

	assert.Equal(t, uint(1), page.Data[0].EpisodeID)
	assert.Equal(t, uint(2), page.Data[1].EpisodeID)

	respWriter = httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodPost, "localhost", nil)

	getEpisodeHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusBadRequest, respWriter.Code)

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

	//passes local. doesnt pass on travis-ci because of gitignore!!
	request = mux.SetURLVars(request, map[string]string{"podcastid": "test", "podcastname": "mypod", "podcastfilename": "sample.mp3"})
	respWriter = httptest.NewRecorder()

	downloadEpisodeHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusOK, respWriter.Code)
	assert.Equal(t, "audio/mpeg", respWriter.Header().Get("Content-Type"))

}

func TestUploadEpisode(t *testing.T) {

	uploadEpisodeHandler := &UploadEpisodeHandler{UserDB: &mocks.MockUserDB{}, PodcastDB: &mocks.MockPodcastDB{}, EpisodeDB: &mocks.MockEpisodeDB{}}

	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)

	episode := models.Episode{
		PodID:   1,
		Created: "stuff",
		Updated: "stuff",
		URL:     "",
		Blurb:   "a blurb",
	}

	// open a testfile
	fileLocation := "../debug_files/test/mypod/sample.mp3"
	file, err := os.Open(fileLocation)
	defer file.Close()

	if err != nil {
		t.Error(err)
	}

	// write file to multipartWriter
	fileWriter, _ := multipartWriter.CreateFormFile("namefile", fileLocation)

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		t.Error(err)
	}

	// add struct to multipart
	fieldWriter, err := multipartWriter.CreateFormField("data")
	if err != nil {
		t.Error(err)
	}

	mepisode, err := json.Marshal(episode)

	if err != nil {
		t.Error(err)
	}

	fieldWriter.Write(mepisode)
	multipartWriter.Close()

	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/upload?podcast=mypodcast", &buf)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	respWriter := httptest.NewRecorder()

	if err != nil {
		t.Error(err)
	}

	uploadEpisodeHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusOK, respWriter.Code)

	// cause fileErr

	buf.Reset()
	multipartWriter = multipart.NewWriter(&buf)
	fieldWriter, err = multipartWriter.CreateFormField("data")
	if err != nil {
		t.Error(err)
	}

	mepisode, err = json.Marshal(episode)

	if err != nil {
		t.Error(err)
	}

	fieldWriter.Write(mepisode)
	multipartWriter.Close()

	respWriter = httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/upload?podcast=mypodcast", &buf)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	uploadEpisodeHandler.ServeHTTP(respWriter, request)
	assert.Equal(t, http.StatusInternalServerError, respWriter.Code)

	// could add full test coverage but ok for now.

}
