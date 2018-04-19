package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

//Categories categoy of podcast
var Categories = []string{"arts", "comedy", "business", "technology", "health", "games", "music", "tv & film",
	"religion & spirituality", "education", "kids & family", "news & politics", "science & medicine", "other", "sports & recreation", "science & medicine"}

// PodcastPage : page data related to podcasts. could of made one struct if i had generics
type PodcastPage struct {
	Data     []Podcast `json:"data"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
}

//EpisodePage page data for Episodes
type EpisodePage struct {
	Data     []Episode `json:"data"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
}

//Config : this is a type used to configure DB, and transmit Singning key.
// Inintially read in from a file
type Config struct {
	Port       string `json:"port"`
	Password   string `json:"password"`
	User       string `json:"user"`
	Schema     string `json:"schema"`
	SigningKey string `json:"signingkey"`
}

//SmtpConfig : configuration for registration email registration confimation
type SmtpConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     string `json:"port"`
}

// ReadFromFile : read to data from a file location
func (s *SmtpConfig) ReadFromFile(location string) {

	file, err := os.Open(location)

	if err != nil {
		log.Println(err)
	}

	err = json.NewDecoder(file).Decode(&s)

	if err != nil {
		log.Println(err)
	}

	fmt.Println(s.Server)
	fmt.Println(s.Username)

}

type TemplateParams struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

//User : Contains a podcast creators admin details incl. email, password etc
type User struct {
	UserName string    `json:"username"`
	Verified bool      `json:"verified"`
	Password string    `json:"password"`
	Token    string    `json:"token"`
	RegToken string    `json:"regtoken"`
	Podcasts []Podcast `json:"podcasts"`
}

//UserTitle : returns username
type UserTitle struct {
	UserName string `json:"username"`
}

//Message : relays a message string to user in a http response
type Message struct {
	Message string `json:"message"`
}

//Podcast : Contains data specific to a podcast incl. name, icon etc
type Podcast struct {
	PodcastID  uint      `gorm:"primary_key"`
	UserEmail  string    `json:"useremail"`
	Icon       string    `json:"icon"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	Downloads  int64     `json:"downloads"`
	Location   string    `json:"location"`
	EpisodeNum int       `json:"episodenum"`
	Details    string    `json:"details"`
	Episodes   []Episode `json:"episodes"`
}

//SecurePodcast : similar to podcast but sent in http response with sensitive user data
type SecurePodcast struct {
	PodcastID  uint      `json:"podcastid"`
	Icon       string    `json:"icon"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	Downloads  int       `json:"downloads"`
	EpisodeNum int       `json:"episodenum"`
	Details    string    `json:"details"`
	Episodes   []Episode `json:"episodes"`
}

//Episode : Contains info on a podcasts episode. a podcast has many episodes
type Episode struct {
	EpisodeID uint   `json:"episodeid"`
	PodID     uint   `json:"podid"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	URL       string `json:"url"`
	Downloads int32  `json:"downloads"`
	Blurb     string `json:"blurb"`
}

//SeedData : read from file used to see debug db
type SeedData struct {
	User     User       `json:"user"`
	Podcasts []*Podcast `json:"podcasts"`
	Episodes []*Episode `json:"episodes"`
}
