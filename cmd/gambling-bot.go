package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	internal "github.com/Namarand/grambling-bot/internal/app"
	"github.com/Ronmi/pastebin"
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

func handleRoll(user twitch.User, contents []string) {
	if !checkPermission(user) {
		return
	}
	if vote == nil {
		sayAdmin("There is no on going vote.")
		return
	}
	if vote.IsOpen {
		sayAdmin("The vote isn't closed.")
		return
	}
	if len(contents) < 3 {
		sayAdmin("You must specify the winner.")
		return
	}
	if !isValidVote(contents[2]) {
		sayAdmin("This is not a correct option: " + contents[2])
	}
	winners := []string{}
	for user, vote := range vote.Vote {
		if vote == contents[2] {
			winners = append(winners, user)
		}
	}
	if len(winners) == 0 {
		sayAdmin("No one vote for this one...")
		return
	}
	sayAdmin("The winner is... " + winners[rand.Intn(len(winners))])
}

func createStat() (string, error) {
	api := pastebin.API{Key: KEY_PASTEBIN}

	transformed := make(map[string][]string)
	for user, value := range vote.Vote {
		transformed[value] = append(transformed[value], user)
	}
	sum := 0
	for _, users := range transformed {
		sum += len(users)
	}
	str := "Total amount: " + strconv.Itoa(sum) + "\n"
	for value, users := range transformed {
		str += value + " (" + strconv.Itoa(len(users)) + "): " + strings.Join(users, ", ") + "\n"
	}
	fmt.Println(str)
	return api.Post(&pastebin.Paste{
		Title:    "Stat Vote",
		Content:  str,
		ExpireAt: pastebin.In1D,
	})
}

func handleDelete(user twitch.User) {
	if !checkPermission(user) {
		return
	}
	vote = nil
}

func handleStat(user twitch.User) {
	if !checkPermission(user) {
		return
	}
	if link, err := createStat(); err == nil {
		sayAdmin(link)
	} else {
		fmt.Println("Error while generating pastebin")
		fmt.Println(err)
	}
}

func handlePrivateStat(user twitch.User) {
	if !checkPermission(user) {
		return
	}
	if link, err := createStat(); err == nil {
		fmt.Println(link)
	} else {
		fmt.Println("Error while generating pastebin")
		fmt.Println(err)
	}
}

func handleReset(user twitch.User) {
	if !checkPermission(user) {
		return
	}
	vote.Vote = make(map[string]string)
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
