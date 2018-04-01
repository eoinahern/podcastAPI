package repository

import (
	"log"

	"github.com/eoinahern/podcastAPI/models"

	"github.com/jinzhu/gorm"
)

type EpisodeDB struct {
	*gorm.DB
}

func (DB *EpisodeDB) GetAllEpisodes(podcastid int) []models.Episode {

	var episodes []models.Episode
	DB.Where("pod_id = ?", podcastid).Find(&episodes)
	return episodes
}

func (DB *EpisodeDB) AddEpisode(episode models.Episode) error {

	db := DB.Save(&episode)

	if db.Error != nil {
		log.Println(db.Error)
	}

	return db.Error

}

func (DB *EpisodeDB) GetLastEpisode() models.Episode {

	var episode models.Episode
	DB.Last(&episode)
	return episode
}
