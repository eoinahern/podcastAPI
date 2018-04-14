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

func TestRegisterHandler(t *testing.T) {

	registerHandler := &RegisterHandler{
		EmailValidator:  &validation.EmailValidation{},
		DB:              &mocks.MockUserDB{},
		MailHelper:      &mocks.MockMailRequest{},
		PassEncryptUtil: &util.PasswordEncryptUtil{},
	}

	responseWriter := httptest.NewRecorder()
	user, _ := json.Marshal(models.User{UserName: "eoin@yahoo.co.uk", Password: "hellothere"})
	request, err := http.NewRequest("POST", "localhost", bytes.NewReader(user))

	if err != nil {
		t.Error(err)
	}

	registerHandler.ServeHTTP(responseWriter, request)
  assert.Equal(t, 200, responseWriter.Code)


}
