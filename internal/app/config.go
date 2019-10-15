package app

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// PastebinCreds is a structure contaning all credentials for Pastebin
type PastebinCreds struct {
	Key string
}

// TwitchCreds is a structure containing all credenials for Twitch
type TwitchCreds struct {
	Channel  string
	Oauth    string
	Username string
}

// Conf is a meta structure containing all nedded configuration for a gambling instance
type Conf struct {
	Pastebin PastebinCreds
	Twitch   TwitchCreds
	Admins   []string
}

// getConf method reads a config file and return and fill a Conf struct
func (c *Conf) getConf(path string) *Conf {

	// Open file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error opening config file : %s", path)
	}
	// Unmarshal it into a Conf struct
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Fatalf("Error reading config file : %s", path)
	}

	return c
}
