package repository

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"github.com/eoinahern/podcastAPI/models"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var columns = []string{"podcast_id", "icon", "name", "category", "downloads", "episode_num", "details"}
var columnsLocation = []string{"podcast_id", "user_email", "icon", "name", "location", "episode_num", "details"}

//setUpMockDB : helper method
func setUpMockDB() (PodcastDB, *sql.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	podcastDB := PodcastDB{DB: db}

	return podcastDB, db, mock

}

func TestGetAllPodcasts(t *testing.T) {
	t.Parallel()

	podcastDB, db, mock := setUpMockDB()
	defer db.Close()

	rows := sqlmock.NewRows(columns).AddRow(1, "icon", "podcast1", "arts", 0, 1, "details about").AddRow(2, "icon.jpeg", "yayrus", "arts", 0, 5, "mo details")
	mock.ExpectQuery("SELECT podcast_id, icon, name").WillReturnRows(rows)

	podcasts := podcastDB.GetAllPodcasts(20, 0, "arts")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, 2, len(podcasts))
	assert.Equal(t, "podcast1", podcasts[0].Name)
	assert.Equal(t, "icon.jpeg", podcasts[1].Icon)

	rows = sqlmock.NewRows(columns)
	mock.ExpectQuery("SELECT podcast_id, icon, name").WillReturnRows(rows)
	podcasts = podcastDB.GetAllPodcasts(50, 20, "")

	assert.Equal(t, 0, len(podcasts))

}

func TestGetPodcast(t *testing.T) {
	t.Parallel()

	podcastDB, db, mock := setUpMockDB()
	defer db.Close()

	username := "me@yahoo.co.uk"
	podcastName := "podcast"

	rows := sqlmock.NewRows(columnsLocation).AddRow(1, username, "icon.jpeg", podcastName, "location", 67, "podcast blurb")
	mock.ExpectQuery("SELECT \\* FROM podcasts WHERE").WithArgs(username, podcastName).WillReturnRows(rows)

	podcast := podcastDB.GetPodcast(username, podcastName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, username, podcast.UserEmail)
	assert.Equal(t, podcastName, podcast.Name)

}

func TestPodcastCreated(t *testing.T) {
	t.Parallel()

	podcastDB, db, mock := setUpMockDB()
	defer db.Close()

	var podID uint = 1
	podcastName := "JRE"

	row := sqlmock.NewRows(columnsLocation).AddRow(podID, "username", "icon", podcastName, "location", 5, "blurb")
	mock.ExpectQuery("SELECT \\* FROM podcasts").WithArgs(podID, podcastName).WillReturnRows(row)

	podcast := podcastDB.CheckPodcastCreated(podID, podcastName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, podID, podcast.PodcastID)
	assert.Equal(t, podcastName, podcast.Name)

}

func TestUpdateNumberPodcasts(t *testing.T) {
	t.Parallel()

	podcastDB, db, mock := setUpMockDB()
	defer db.Close()

	rows := sqlmock.NewRows(columnsLocation).AddRow(1, "email", "", "podcast", "location", 0, "a podcast")
	mock.ExpectQuery("SELECT \\* FROM podcasts").WithArgs(1).WillReturnRows(rows)
	mock.ExpectPrepare("UPDATE podcasts SET episode_num")
	mock.ExpectExec("UPDATE podcasts SET episode_num").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	podcastDB.UpdatePodcastNumEpisodes(1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

}

func TestCreatePodcast(t *testing.T) {
	t.Parallel()

	podcastDB, db, mock := setUpMockDB()
	defer db.Close()

	podcast := &models.Podcast{UserEmail: "eoin", Icon: "none", Name: "name", Category: "arts", Downloads: 0, Location: "location/location", Details: "podcast about something"}
	mock.ExpectPrepare("INSERT INTO podcasts")
	mock.ExpectExec("INSERT INTO podcasts").WithArgs("eoin", "none", "name", "arts", 0, "location/location", "podcast about something").WillReturnResult(sqlmock.NewResult(1, 1))

	errorCreate := podcastDB.CreatePodcast(podcast)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, nil, errorCreate)

}
