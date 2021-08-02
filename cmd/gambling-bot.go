package main

import (
	"os"

	"github.com/apex/log"
	loghandler "github.com/apex/log/handlers/logfmt"

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
	app.Version = "1.2.4"

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

		// logs
		logsSetup()

		// Create a new gambling instance
		gambling := internal.NewGambling(c.String("config"))

		// Start it
		return gambling.Start()

	}

	// Run
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}

}

// logsSetup is used to set up logs from environment variable
func logsSetup() {

	// choose log type
	log.SetHandler(loghandler.Default)
	// loglevel
	loglvl := os.Getenv("LOGLEVEL")
	if loglvl == "" {
		log.SetLevelFromString("Warn")
		log.Warn("No log level environment variable configured, defaulted to WARN")
	} else {
		log.SetLevelFromString(loglvl)
		log.Info("Log level set")
	}

}
