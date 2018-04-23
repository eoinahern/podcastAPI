package util

import (
	"fmt"

	"github.com/eoinahern/podcastAPI/models"
)

const limitStr string = "limit="
const offsetStr string = "offset="

// CreatePodcastPage estimating how ill construct the page object to return??
func CreatePodcastPage(endpoint string, limit int, offset int, totalItems int) *models.PodcastPage {

	return &models.PodcastPage{Data: []models.Podcast{},
		Next:     createNextURL(endpoint, limit, offset, totalItems),
		Previous: createPreviousURL(endpoint, limit, offset)}
}

//CreateEpisodePage guesstimate TODO
func CreateEpisodePage(endpoint string, limit int, offset int, totalItems int) *models.EpisodePage {

	return &models.EpisodePage{Data: []models.Episode{},
		Next:     createNextURL(endpoint, limit, offset, totalItems),
		Previous: createPreviousURL(endpoint, limit, offset)}
}

func createNextURL(endpoint string, limit int, offset int, totalItems int) string {

	//we need total in db here???

	var result string

	if (offset + limit) >= totalItems {
		return result
	}

	result = fmt.Sprintf("%s?%s%d&%s%d", endpoint, limitStr, limit, offsetStr, offset+limit)

	return result
}

func createPreviousURL(endpoint string, limit int, offset int) string {

	var result string

	if offset == 0 || offset-limit < 0 {
		return result
	}

	result = fmt.Sprintf("%s?%s%d&%s%d", endpoint, limitStr, limit, offsetStr, offset-limit)
	return result

}
