package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

	//call with empty params
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

}
