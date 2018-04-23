package mocks

import (
	"strings"

	"github.com/eoinahern/podcastAPI/models"
)

//MockUserDB : mocks my userDB wrapper
type MockUserDB struct {
}

//CountRows : mocking count rows
func (DB *MockUserDB) CountRows() int {
	return 1
}

//CheckExist check exist mock
func (DB *MockUserDB) CheckExist(email string) bool {
	return false
}

//ValidateUserPlusRegToken mock imp
func (DB *MockUserDB) ValidateUserPlusRegToken(email string, regToken string) bool {

	if strings.Compare(email, "eoin@yahoo.com") == 0 && strings.Compare(regToken, "1234") == 0 {
		return true
	}

	return false
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

func (DB *MockPodcastDB) GetAll(limit uint16, offset uint16, by string) []models.SecurePodcast {
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

func (DB *MockPodcastDB) GetPodcast(username string, podcastname string) *models.Podcast {
	return &models.Podcast{
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

func (DB *MockPodcastDB) CreatePodcast(podcast *models.Podcast) error {
	return nil
}

func (DB *MockPodcastDB) UpdatePodcastNumEpisodes(id uint) {

}

type MockEpisodeDB struct {
}

func getMockEpisode() models.Episode {
	return models.Episode{
		EpisodeID: 1,
		PodID:     2,
		Created:   "",
		Updated:   "",
		URL:       "google.com/files/episode",
		Downloads: 200,
		Blurb:     "episode featuring bill oherlihy",
	}
}

func (DB *MockEpisodeDB) CountRows() int {
	return 1
}

func CountRowsByID(podID int) int {
	return 90
}

func (DB *MockEpisodeDB) GetAllEpisodes(podcastid int, limit uint16, offset uint16) []models.Episode {

	return []models.Episode{
		{
			EpisodeID: 1,
			PodID:     2,
			Created:   "",
			Updated:   "",
			URL:       "google.com/files/episode",
			Downloads: 200,
			Blurb:     "episode featuring bill oherlihy",
		},
		{
			EpisodeID: 2,
			PodID:     3,
			Created:   "",
			Updated:   "",
			URL:       "google.com/files/episode2",
			Downloads: 117,
			Blurb:     "episode featuring lee syatt",
		},
	}

}

func (DB *MockEpisodeDB) AddEpisode(episode models.Episode) error {
	return nil
}

func (DB *MockEpisodeDB) GetSingleEpisode(podcastID uint, episodeID uint) models.Episode {
	return getMockEpisode()
}

func (DB *MockEpisodeDB) GetLastEpisode() models.Episode {
	return getMockEpisode()
}
