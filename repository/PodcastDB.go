package repository

import (
	"database/sql"
	"log"

	"github.com/eoinahern/podcastAPI/models"
)

//PodcastDBInt interface
type PodcastDBInt interface {
	CountRows() int
	GetAll(limit uint16, offset uint16, by string) []models.SecurePodcast
	GetPodcast(username string, podcastname string) *models.Podcast
	CheckPodcastCreated(podID uint, podname string) models.Podcast
	CreatePodcast(podcast *models.Podcast) error
	UpdatePodcastNumEpisodes(id uint)
}

//PodcastDB : podcast database helper
type PodcastDB struct {
	*sql.DB
}

//CountRows : num rows
func (DB *PodcastDB) CountRows() int {

	var count int
	row := DB.QueryRow("SELECT COUNT(*) FROM podcasts")
	err := row.Scan(&count)

	if err == nil {
		return count
	}

	return 0
}

//GetAll : get all podcasts. not episodes just a podcast name!!
//TODO: need to page. potentially filter by category etc here!! popularity etc
func (DB *PodcastDB) GetAll(limit uint16, offset uint16, by string) []models.SecurePodcast {

	var podcasts []models.SecurePodcast

	rows, err := DB.Query("SELECT podcast_id, icon, name, episode_num, details from podcasts")
	defer rows.Close()

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {

		var securePodcast models.SecurePodcast
		if err := rows.Scan(&securePodcast.PodcastID, &securePodcast.Icon,
			&securePodcast.Name, &securePodcast.EpisodeNum, &securePodcast.Details); err != nil {

			log.Println(err)

		} else {

			podcasts = append(podcasts, securePodcast)
		}

	}

	return podcasts
}

//GetPodcast : get a podcast from the DB based on username and podcastName
//probably more for admin use as have to pass email?
func (DB *PodcastDB) GetPodcast(userName string, podcastName string) *models.Podcast {

	var podcast models.Podcast
	row := DB.QueryRow("SELECT * FROM podcasts WHERE user_email = ? AND name = ?", userName, podcastName)

	err := row.Scan(&podcast.PodcastID, &podcast.UserEmail, &podcast.Icon, &podcast.Name,
		&podcast.Location, &podcast.EpisodeNum, &podcast.Details)

	if err != nil {
		log.Println(err)
	}

	return &podcast

}

//CheckPodcastCreated : check if this podcast exists in DB
func (DB *PodcastDB) CheckPodcastCreated(podcastID uint, podcastName string) models.Podcast {

	var podcast models.Podcast
	row := DB.QueryRow("SELECT * FROM podcasts WHERE name = ? AND podcast_id = ?", podcastID, podcastName)
	err := row.Scan(&podcast.PodcastID, &podcast.UserEmail, &podcast.Icon, &podcast.Name,
		&podcast.Location, &podcast.EpisodeNum, &podcast.Details)

	if err != nil {
		log.Println(err)
	}

	return podcast

}

//UpdatePodcastNumEpisodes : update number of episodes
func (DB *PodcastDB) UpdatePodcastNumEpisodes(id uint) {

	var podcast models.Podcast

	row := DB.QueryRow("SELECT * FROM podcasts WHERE podcast_id = ?", id)
	row.Scan(&podcast.PodcastID, &podcast.UserEmail, &podcast.Icon, &podcast.Name,
		&podcast.Location, &podcast.EpisodeNum, &podcast.Details)

	podcast.EpisodeNum++

	stmt, err := DB.Prepare("UPDATE podcasts SET episode_num = ? WHERE podcast_id= ?")
	defer stmt.Close()

	if err != nil {
		log.Println("problem with stmt")
		log.Println(err)
	}

	_, err = stmt.Exec(podcast.EpisodeNum, id)

	if err != nil {
		log.Println(err)
	}

}

//CreatePodcast : save podcast to database
func (DB *PodcastDB) CreatePodcast(podcast *models.Podcast) error {

	stmt, err := DB.Prepare("INSERT INTO podcasts(user_email, icon, name, category, downloads, location, details) VALUES(?,?,?,?,?,?,?)")

	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(podcast.UserEmail, podcast.Icon, podcast.Name, podcast.Category, podcast.Downloads, podcast.Location, podcast.Details)

	return err
}
