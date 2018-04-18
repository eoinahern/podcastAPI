package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/eoinahern/podcastAPI/models"
)

//EpisodeDBInt interface
type EpisodeDBInt interface {
	CountRows() int
	GetAllEpisodes(podcastid int) []models.Episode
	AddEpisode(episode models.Episode) error
	GetSingleEpisode(podcastID uint, episodeID uint) models.Episode
	GetLastEpisode() models.Episode
}

//EpisodeDB : collect, maintain epoisode data in DB
type EpisodeDB struct {
	*sql.DB
}

//CountRows : num rows
func (DB *EpisodeDB) CountRows() int {

	var count int
	row := DB.QueryRow("SELECT COUNT(*) FROM episodes")
	err := row.Scan(&count)

	if err == nil {
		return count
	}

	return 0
}

//GetAllEpisodes : get all episodes associated with specific podcast
func (DB *EpisodeDB) GetAllEpisodes(podcastid int) []models.Episode {

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

}

//AddEpisode : Add episode data to database
func (DB *EpisodeDB) AddEpisode(episode models.Episode) error {

	stmt, err := DB.Prepare("INSERT INTO episodes(pod_id, created, updated, url, downloads, blurb) VALUES(?,?,?,?,?,?)")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
	}

	res, err := stmt.Exec(episode.PodID, episode.Created, episode.Updated, episode.URL, episode.Downloads, episode.Blurb)

	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return err
	}

	insertID, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println(fmt.Sprintf("last insert id :  %d ", insertID))
	fmt.Println(fmt.Sprintf("rows affected:  %d ", rowsAffected))

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
	row, err := DB.Query("SELECT TOP 1 * FROM episodes ORDER BY episode_id DESC LIMIT 1") //DESC wasnt working on QueryRow() ????
	defer row.Close()

	if err != nil {
		log.Println(err)
	}

	for row.Next() {
		err = row.Scan(&episode.EpisodeID, &episode.PodID, &episode.Created, &episode.Updated, &episode.URL, &episode.Downloads, &episode.Blurb)

		if err != nil {
			log.Println(err)
		}

	}

	return episode
}
