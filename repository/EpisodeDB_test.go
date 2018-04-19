package repository

import (
	"database/sql"
	"testing"

	"github.com/eoinahern/podcastAPI/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var episodeColumns = []string{"episode_id", "pod_id", "created", "updated", "url", "downloads", "blurb"}
var episode1 = &models.Episode{EpisodeID: 1, PodID: 1, Created: "", Updated: "", URL: "files/podcasts", Downloads: 1, Blurb: "podcast about stuff"}
var episode2 = &models.Episode{EpisodeID: 2, PodID: 1, Created: "", Updated: "", URL: "files/podcasts", Downloads: 200, Blurb: "podcast about fishing"}

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

	episodeDB, db, mock := setUpEpisodeMocks()

	defer db.Close()

	rows := sqlmock.NewRows(episodeColumns).AddRow(episode1.EpisodeID, episode1.PodID, episode1.Created, episode1.Updated, episode1.URL, episode1.Downloads, episode1.Blurb)
	rows = rows.AddRow(episode2.EpisodeID, episode2.PodID, episode2.Created, episode2.Updated, episode2.URL, episode2.Downloads, episode2.Blurb)
	mock.ExpectQuery("SELECT \\* FROM episodes").WithArgs(episode1.PodID).WillReturnRows(rows)

	episodes := episodeDB.GetAllEpisodes(int(episode1.PodID), 20, 0)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, 2, len(episodes))
	assert.Equal(t, 1, int(episodes[0].PodID))
	assert.Equal(t, 1, int(episodes[1].PodID))

}

func TestAddEpisode(t *testing.T) {

	t.Parallel()

	episodeDB, db, mock := setUpEpisodeMocks()
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO episodes")
	mock.ExpectExec("INSERT INTO episodes").WithArgs(episode1.PodID, episode1.Created, episode1.Updated, episode1.URL, episode1.Downloads, episode1.Blurb).WillReturnResult(sqlmock.NewResult(1, 1))

	err := episodeDB.AddEpisode(*episode1)

	if errExpt := mock.ExpectationsWereMet(); errExpt != nil {
		t.Errorf("err %s", errExpt)
	}

	assert.Equal(t, nil, err)

}

func TestGetSingleEpisode(t *testing.T) {

	t.Parallel()
	episodeDB, db, mock := setUpEpisodeMocks()
	defer db.Close()

	rows := sqlmock.NewRows(episodeColumns).AddRow(episode1.EpisodeID, episode1.PodID, episode1.Created, episode1.Updated, episode1.URL, episode1.Downloads, episode1.Blurb)
	mock.ExpectQuery("SELECT").WithArgs(episode1.EpisodeID, episode1.PodID).WillReturnRows(rows)

	episode := episodeDB.GetSingleEpisode(episode1.EpisodeID, episode1.PodID)

	if errExpt := mock.ExpectationsWereMet(); errExpt != nil {
		t.Errorf("err %s", errExpt)
	}

	assert.Equal(t, 1, int(episode.PodID))
	assert.Equal(t, episode1.URL, episode.URL)

}

func TestGetLastEpisode(t *testing.T) {

	t.Parallel()

	episodeDB, db, mock := setUpEpisodeMocks()
	defer db.Close()

	rows := sqlmock.NewRows(episodeColumns).AddRow(episode1.EpisodeID, episode1.PodID, episode1.Created, episode1.Updated, episode1.URL, episode1.Downloads, episode1.Blurb).AddRow(episode2.EpisodeID, episode2.PodID, episode2.Created, episode2.Updated, episode2.URL, episode2.Downloads, episode2.Blurb)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	episode := episodeDB.GetLastEpisode()

	if errExpt := mock.ExpectationsWereMet(); errExpt != nil {
		t.Errorf("err %s", errExpt)
	}

	assert.Equal(t, episode2.EpisodeID, episode.EpisodeID)

}
