package util

import (
	"fmt"
	"net/http"

	"github.com/eoinahern/podcastAPI/models"
)

const limitStr string = "limit="
const offsetStr string = "offset="
const podIDStr string = "pod_id="

// CreatePodcastPage estimating how ill construct the page object to return??
func CreatePodcastPage(endpoint *http.Request, limit int, offset int, totalItems int) *models.PodcastPage {

	return &models.PodcastPage{Data: []models.Podcast{},
		Next:     createNextURL(endpoint, 0, limit, offset, totalItems),
		Previous: createPreviousURL(endpoint, 0, limit, offset)}
}

// CreateEpisodePage used to create page data struct related to episodes
func CreateEpisodePage(endpoint *http.Request, podid int, limit int, offset int, totalItems int) *models.EpisodePage {

	return &models.EpisodePage{Data: []models.Episode{},
		Next:     createNextURL(endpoint, podid, limit, offset, totalItems),
		Previous: createPreviousURL(endpoint, podid, limit, offset)}
}

func createNextURL(endpoint *http.Request, podid int, limit int, offset int, totalItems int) string {

	if (offset + limit) >= totalItems {
		return ""
	}

	fmt.Println(totalItems)

	return createURL(endpoint, podid, limit, offset+limit)
}

func createURL(endpoint *http.Request, podid int, limit int, newOffset int) string {

	var result string

	if podid == 0 {
		result = fmt.Sprintf("%s%s?%s%d&%s%d", endpoint.URL.Host, endpoint.URL.Path, limitStr, limit, offsetStr, newOffset)
	} else {
		result = fmt.Sprintf("%s%s?%s%d&%s%d&%s%d", endpoint.URL.Host, endpoint.URL.Path, podIDStr, podid, limitStr, limit, offsetStr, newOffset)
	}

	return result
}

func createPreviousURL(endpoint *http.Request, podid int, limit int, offset int) string {

	if offset == 0 || offset-limit < 0 {
		return ""
	}

	return createURL(endpoint, podid, limit, offset-limit)

}
