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

func (DB *MockPodcastDB) CountRows() int {
	return 1
}

func (DB *MockPodcastDB) GetAll() []models.SecurePodcast {
	return getTestPodcasts()
}

func getTestPodcasts() []models.SecurePodcast {
	return []models.SecurePodcast{
		{
			PodcastID:  1,
			Icon:       "",
			Name:       "church of whats happening now",
			EpisodeNum: 2,
			Details:    "edibles and the christ killer",
		},
		{
			PodcastID:  2,
			Icon:       "",
			Name:       "jre",
			EpisodeNum: 2,
			Details:    " a podcast about stuff",
		},
	}
}

func (DB *MockPodcastDB) GetPodcast(username string, podcastname string) models.Podcast {
	return models.Podcast{
		PodcastID:  1,
		Icon:       "",
		Name:       podcastname,
		EpisodeNum: 2,
		Details:    "edibles and the christ killer",
	}

}

func (DB *MockPodcastDB) CheckPodcastCreated(podID uint, podname string) models.Podcast {

	return models.Podcast{
		PodcastID:  podID,
		Icon:       "",
		Name:       podname,
		EpisodeNum: 2,
		Details:    "edibles and the christ killer",
	}

}

func (DB *MockPodcastDB) CreatePodcast(podcast *models.Podcast) {

}

func (DB *MockPodcastDB) UpdatePodcastNumEpisodes(id uint) {

}
