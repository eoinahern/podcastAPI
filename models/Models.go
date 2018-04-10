package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

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
	UserName string    `json:"username" gorm:"primary_key; type:VARCHAR(80)"`
	Verified bool      `json:"verified" gorm:"type:BOOLEAN; default:false" `
	Password string    `json:"password" gorm:"type:TEXT"`
	Token    string    `json:"token" sql:"-" gorm:"-"`
	RegToken string    `json:"regtoken" gorm:"type:TEXT"`
	Podcasts []Podcast `json:"podcasts" gorm:"ForeignKey:UserEmail"`
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
	PodcastID  uint      `gorm:"primary_key"  json:"podcastid"`
	UserEmail  string    `json:"useremail" gorm:"type:VARCHAR(80)"`
	Icon       string    `json:"icon"   gorm:"type:TEXT"`
	Name       string    `json:"name" gorm:"type: TEXT"`
	Location   string    `json:"location" gorm:"type:VARCHAR(100)"`
	EpisodeNum int       `json:"episodenum" gorm:"type:INTEGER; default:0"`
	Details    string    `json:"details" gorm:"type:TEXT"`
	Episodes   []Episode `json:"episodes"  gorm:"ForeignKey:PodID"`
}

//SecurePodcast : similar to podcast but sent in http response with sensitive user data
type SecurePodcast struct {
	PodcastID  uint      `json:"podcastid"`
	Icon       string    `json:"icon"`
	Name       string    `json:"name" `
	EpisodeNum int       `json:"episodenum"`
	Details    string    `json:"details"`
	Episodes   []Episode `json:"episodes"`
}

//Episode : Contains info on a podcasts episode. a podcast has many episodes
type Episode struct {
	EpisodeID uint   `json:"episodeid" gorm:"primary_key"`
	PodID     uint   `json:"podid"`
	Created   string `json:"created" gorm:"type: TEXT"`
	Updated   string `json:"updated" gorm:"type: TEXT"`
	URL       string `json:"url" gorm:"type: TEXT"`
	Downloads int32  `json:"downloads" gorm:"type:INTEGER; not null default:0"`
	Blurb     string `json:"blurb" gorm:"type: TEXT"`
}

//SeedData : read from file used to see debug db
type SeedData struct {
	User     User      `json:"user"`
	Podcasts []Podcast `json:"podcasts"`
	Episodes []Episode `json:"episodes"`
}
