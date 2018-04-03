package repository

import (
	"database/sql"

	"github.com/eoinahern/podcastAPI/models"
)

//EpisodeDB : collect, maintain epoisode data in DB
type EpisodeDB struct {
	*sql.DB
}

//GetAllEpisodes : get all episodes associated with specific podcast
func (DB *EpisodeDB) GetAllEpisodes(podcastid int) []models.Episode {

	/*var episodes []models.Episode
	DB.Where("pod_id = ?", podcastid).Find(&episodes)
	return episodes*/
	return nil
}

//AddEpisode : Add episode data to database
func (DB *EpisodeDB) AddEpisode(episode models.Episode) error {

	/*db := DB.Save(&episode)

	if db.Error != nil {
		log.Println(db.Error)
	}

	return db.Error*/
	return nil

}

//GetSingleEpisode : get data about episode base on id, and podname. maybe id aswell.
func (DB *EpisodeDB) GetSingleEpisode(podcastid string, podcastname string) {

}

//GetLastEpisode : get last episode from db? not sure if this is required?
func (DB *EpisodeDB) GetLastEpisode() models.Episode {

	/*var episode models.Episode
	DB.Last(&episode)
	return episode*/return models.Episode{}
}
