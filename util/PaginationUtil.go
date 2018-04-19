package util

import (
	"fmt"

	"github.com/eoinahern/podcastAPI/models"
)

const endpoint string = "http:localhost/8080/"
const limitStr string = "limit="
const offsetStr string = "offset="

// CreatePodcastPage estimating how ill construct the page object to return??
func CreatePodcastPage(limit uint16, offset uint16) models.PodcastPage {

	return models.PodcastPage{Data: []models.Podcast{},
		Next:     createNextURL("podcast/", limit, offset),
		Previous: createPreviousURL("podcast/", limit, offset)}
}

//CreateEpisodePage guesstimate TODO
func CreateEpisodePage(limit uint16, offset uint16) models.EpisodePage {

	return models.EpisodePage{Data: []models.Episode{},
		Next:     createNextURL("podcast/", limit, offset),
		Previous: createPreviousURL("podcast/", limit, offset)}
}

func createNextURL(resource string, limit uint16, offset uint16) string {

	//we need total in db here????
	var result string
	result = fmt.Sprintf("%s%s?%s%d", endpoint, resource, limitStr, limit)

	return result
}

func createPreviousURL(resource string, limit uint16, offset uint16) string {

	var result string

	if offset == 0 {
		return result
	}

	result = fmt.Sprintf("%s%s?%s%d", endpoint, resource, limitStr, limit)
	return result

}
