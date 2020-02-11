package main

import (
	"log"
	"os"

	internal "github.com/Namarand/gambling-bot/internal/app"
	"github.com/urfave/cli"
)

// Entrypoint
func main() {

	// Declare app
	app := cli.NewApp()
	// Basic config
	app.Name = "gambling-bot"
	app.Usage = "A golang powered Twitch bot handling simple vote mechanism"
	app.Version = "0.2.1"

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
