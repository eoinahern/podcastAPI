package repository

import (
	"database/sql"
	"testing"

	"github.com/eoinahern/podcastAPI/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var episodeColumns = []string{"episode_id", "pod_id", "created", "updated", "url", "downloads", "blurb"}
var episode1 = &models.Episode{EpisodeID: 1, PodID: 1, Created: "", Updated: "", URL: "files/podcasts", Downloads: 1, Blurb: "podcast about stuff"}
var episode2 = &models.Episode{EpisodeID: 2, PodID: 3, Created: "", Updated: "", URL: "files/podcasts", Downloads: 200, Blurb: "podcast about fishing"}

//setUpEpisodeMocks : create mocks helper
func setUpEpisodeMocks() (EpisodeDB, *sql.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	episodeDB := EpisodeDB{DB: db}
	return episodeDB, db, mock
}

func TestGetAllEPisodes(t *testing.T) {

	t.Parallel()

}

func TestAddEpisode(t *testing.T) {

	t.Parallel()

}

func TestGetSingleEpisode(t *testing.T) {

	t.Parallel()

}

func TestGetLastEpisode(t *testing.T) {

}
