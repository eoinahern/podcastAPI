package repository

import (
	"database/sql"
	"log"

	"github.com/eoinahern/podcastAPI/models"
)

//EpisodeDB : collect, maintain epoisode data in DB
type EpisodeDB struct {
	*sql.DB
}

//GetAllEpisodes : get all episodes associated with specific podcast
func (DB *EpisodeDB) GetAllEpisodes(podcastid uint) []models.Episode {

	var episodes []models.Episode
	rows, err := DB.Query("SELECT * FROM episodes WHERE pod_id = ?", podcastid)

	defer rows.Close()

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var episode models.Episode
		err = rows.Scan(&episode.EpisodeID, &episode.PodID, &episode.Created, &episode.Updated, &episode.URL, &episode.Downloads, &episode.Blurb)

		if err != nil {
			log.Println(err)
		} else {
			episodes = append(episodes, episode)
		}
	}

	return episodes

	/*var episodes []models.Episode
	DB.Where("pod_id = ?", podcastid).Find(&episodes)
	return episodes*/

}

//AddEpisode : Add episode data to database
func (DB *EpisodeDB) AddEpisode(episode models.Episode) error {

	stmt, err := DB.Prepare("INSERT INTO episodes(episode_id, pod_id, created, updated, url, downloads, blurb) VALUES(?,?,?,?,?,?,?)")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
	}

	_, err = stmt.Exec(episode.EpisodeID, episode.PodID, episode.Created, episode.Updated, episode.URL, episode.Downloads, episode.Blurb)
	return err

}

//GetSingleEpisode : get data about episode base on id, and podname. maybe id aswell.
func (DB *EpisodeDB) GetSingleEpisode(podcastid uint, episodeID uint) models.Episode {

	var episode models.Episode
	row := DB.QueryRow("SELECT * FROM episodes WHERE episode_id = ? AND pod_id = ?", episodeID, podcastid)
	err := row.Scan(&episode.EpisodeID, &episode.PodID, &episode.Created, &episode.Updated, &episode.URL, &episode.Downloads, &episode.Blurb)

	if err != nil {
		log.Println(err)
	}

	return episode
}

//GetLastEpisode : get last episode from db? not sure if this is required?
func (DB *EpisodeDB) GetLastEpisode() models.Episode {

	//SELECT * FROM TABLE ORDER BY episode_id DESC LIMIT 1
	var episode models.Episode
	row := DB.QueryRow("SELECT * FROM episodes ORDER BY episode_id DESC LIMIT 1")
	err := row.Scan(&episode.EpisodeID, &episode.PodID, &episode.Created, &episode.Updated, &episode.URL, &episode.Downloads, &episode.Blurb)

	if err != nil {
		log.Println(err)
	}

	return episode
}
