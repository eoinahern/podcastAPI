package repository

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"github.com/eoinahern/podcastAPI/models"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetAll(t *testing.T) {

}

func TestGetPodcast(t *testing.T) {

}

func TestPodcastCreated(t *testing.T) {

}

func TestUpdateNumberPodcasts(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	podcastDB := PodcastDB{DB: db}
	rows := sqlmock.NewRows([]string{"podcast_id", "user_email", "icon", "name", "location", "episode_num", "details"}).AddRow(1, "email", "", "podcast", "location", 0, "a podcast")
	mock.ExpectQuery("SELECT \\* FROM podcasts").WithArgs(1).WillReturnRows(rows)
	mock.ExpectPrepare("UPDATE podcasts SET episode_num")
	mock.ExpectExec("UPDATE podcasts SET episode_num").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	podcastDB.UpdatePodcastNumEpisodes(1)

}

func TestCreatePodcast(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	podcastDB := PodcastDB{DB: db}

	podcast := &models.Podcast{UserEmail: "eoin", Icon: "none", Name: "name", Location: "location/location", Details: "podcast about something"}
	mock.ExpectPrepare("INSERT INTO podcasts")
	mock.ExpectExec("INSERT INTO podcasts").WithArgs("eoin", "none", "name", "location/location", "podcast about something").WillReturnResult(sqlmock.NewResult(1, 1))

	errorCreate := podcastDB.CreatePodcast(podcast)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, nil, errorCreate)

}
