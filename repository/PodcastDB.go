package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/eoinahern/podcastAPI/models"
)

type PodcastDB struct {
	*sql.DB
}

func (DB *PodcastDB) GetAll() []models.SecurePodcast {

	/*	var podcasts []models.SecurePodcast
		rows, err := DB.Raw("SELECT podcast_id, icon, name, episode_num from podcasts").Rows()

		if err != nil {
			log.Println(err)
		}

		defer rows.Close()

		for rows.Next() {
			var pod models.SecurePodcast
			rows.Scan(&pod.PodcastID, &pod.Icon, &pod.Name, &pod.EpisodeNum)
			podcasts = append(podcasts, pod)
		}

		return podcasts*/
	return []models.SecurePodcast{}
}

func (DB *PodcastDB) GetPodcast(userName string, podcastName string) *models.Podcast {
	/*var podcast models.Podcast
	DB.Where("user_email = ? AND name = ?", userName, podcastName).First(&podcast)
	return &podcast*/

	return &models.Podcast{}
}

func (DB *PodcastDB) CheckPodcastCreated(podcastID uint, podcastName string) models.Podcast {

	/*var podcast models.Podcast
	DB.Where("name = ? AND  podcast_id = ?", podcastName, podcastID).First(&podcast)

	return podcast*/
	return models.Podcast{}

}

func (DB *PodcastDB) UpdatePodcastNumEpisodes(id uint) {

	/*var podcast models.Podcast
	DB.Where("podcast_id = ?", id).First(&podcast)
	podcast.EpisodeNum += 1
	db := DB.Save(&podcast)

	if db.Error != nil {
		log.Println(db.Error)
	}*/

}

//CreatePodcast : save podcast to database
func (DB *PodcastDB) CreatePodcast(podcast *models.Podcast) error {

	stmt, err := DB.Prepare("INSERT INTO podcasts(user_email, icon, name, location, details) VALUES(?,?,?,?,?)")

	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()
	res, err := stmt.Exec(podcast.UserEmail, podcast.Icon, podcast.Name, podcast.Location, podcast.Details)

	rows, _ := res.RowsAffected()
	fmt.Println(fmt.Sprintf("num rows affected %s", string(rows)))

	return err
}
