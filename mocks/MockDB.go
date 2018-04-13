package mocks

import (
	"github.com/eoinahern/podcastAPI/models"
)

//MockUserDB : mocks my userDB wrapper
type MockUserDB struct {
}

//CountRows : mocking count rows
func CountRows() int {
	return 1
}

//CheckExist check exist mock
func (DB *MockUserDB) CheckExist(email string) bool {
	return true
}

//ValidateUserPlusRegToken mock imp
func (DB *MockUserDB) ValidateUserPlusRegToken(email string, regToken string) bool {
	return true
}

//SetVerified mock setVerified
func (DB *MockUserDB) SetVerified(email string, regToken string) {

}

//ValidatePasswordAndUser validate pass returns tur regardless of input
func (DB *MockUserDB) ValidatePasswordAndUser(email string, password string) bool {
	return true
}

//Insert mock insert does nothing
func (DB *MockUserDB) Insert(user *models.User) {

}

//GetUser mockuser return
func (DB *MockUserDB) GetUser(email string) models.User {

	return models.User{UserName: "jimmy@yahoo.com",
		Verified: true,
		Password: "pass",
		Token:    "token",
		RegToken: "reg"}
}

//MockPodcastDB mock podcastDB wrapper
type MockPodcastDB struct {
}
