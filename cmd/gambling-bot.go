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

func handleCreate(user twitch.User, contents []string) {
	if !checkPermission(user) {
		return
	}
	if vote != nil {
		sayAdmin("There is already a vote going, you should delete it first with '!gamble delete'.")
		return
	}
	if len(contents) <= 2 {
		sayAdmin("You must pass the possibilities as arguments.")
		return
	}
	vote = new(Gambling)
	vote.IsOpen = true
	vote.Vote = make(map[string]string)
	vote.Possibilities = contents[2:]
	sayAdmin("There is a new vote! You can vote with '!gamble vote <vote>'.")
	sayAdmin("You can vote for: " + strings.Join(contents[2:], ", "))
}

func handleClose(user twitch.User) {
	if !checkPermission(user) {
		return
	}
	if vote == nil {
		sayAdmin("There is no on going vote.")
		return
	}
	vote.IsOpen = false
	sayAdmin("The vote is now closed!")
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

func handleVote(user twitch.User, contents []string) {
	if vote == nil || !vote.IsOpen {
		fmt.Println("Don't try to vote while it's close...")
		return
	}
	if len(contents) < 2 {
		fmt.Println("Expected one argument to vote...")
	}
	if isValidVote(contents[2]) {
		vote.Vote[user.Name] = contents[2]
	} else {
		fmt.Println(user.Name + ": Invalid vote: got '" + contents[2] + "'")
	}
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

func parseMessage(message twitch.PrivateMessage) {
	if !strings.HasPrefix(message.Message, "!gamble") {
		return
	}
	contents := strings.Split(message.Message, " ")
	if len(contents) == 1 {
		return
	}
	switch contents[1] {
	case "create":
		handleCreate(message.User, contents)
	case "close":
		handleClose(message.User)
	case "roll":
		handleRoll(message.User, contents)
	case "vote":
		handleVote(message.User, contents)
	case "delete":
		handleDelete(message.User)
	case "reset":
		handleReset(message.User)
	case "stats":
		handleStat(message.User)
	case "privatestats":
		handlePrivateStat(message.User)
	}
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
			Value:  "/etc/gramble/config.yml",
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
