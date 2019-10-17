package main

import (
	"log"
	"os"

	internal "github.com/Namarand/grambling-bot/internal/app"
	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/urfave/cli"
)

type Gambling struct {
	IsOpen        bool
	Possibilities []string
	Vote          map[string]string
}

var vote *Gambling
var client *twitch.Client
var CHANNEL = "val_pl_magicarenafr"
var KEY_PASTEBIN = "58a788a403a74613ed74e745f473aaa6"

func isValidVote(str string) bool {
	if vote == nil {
		return false
	}
	for _, key := range vote.Possibilities {
		if key == str {
			return true
		}
	}
	return false
}

func sayAdmin(message string) {
	client.Say(CHANNEL, "@"+CHANNEL+": "+message)
}

func checkPermission(user twitch.User) bool {
	return user.Name == "namarand" || user.Name == CHANNEL
}

func main() {

	// Declare app
	app := cli.NewApp()
	// Basic config
	app.Name = "gambling-bot"
	app.Usage = "A golang powered Twitch bot handling simple vote mechanism"
	app.Version = "0.1.0"

	// Flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Path to yaml config file",
			EnvVar: "CONFIG",
			Value:  "/etc/gamble/config.yml",
		},
	}

	// Action
	app.Action = func(c *cli.Context) error {

		// Create a new gambling instance
		gambling := internal.NewGambling(c.String("config"))

		// Start it
		return gambling.Start()

	}

	// Run
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
