package app

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

// Vote is a structure handling all voting params and status
type Vote struct {
	IsOpen        bool
	Possibilities []string
	Votes         map[string]string
}

// Gambling is a meta structure containing all the stuff needed by a Gambling instance
type Gambling struct {
	// Config from yaml file
	Config Conf
	// Twitch client
	Twitch *twitch.Client
	// Vote
	CurrentVote *Vote
}

// NewGambling func create a new Gambling struct
func NewGambling(confPath string) *Gambling {
	// Empty new struct
	g := new(Gambling)

	// Parse config
	g.Config.getConf(confPath)

	// Setup twitch client
	g.Twitch = twitch.NewClient(g.Config.Twitch.Username, g.Config.Twitch.Oauth)

	// Plug function on Twitch events
	g.twitchOnEventSetup()

	// Join Twitch channel
	g.join()

	// init vote
	g.CurrentVote = new(Vote)

	return g

}

// Extract cmd and args from message
func extractCommand(message string) (string, []string) {

	contents := strings.Fields(message)

	if len(contents) == 1 {
		return "", nil
	}

	if len(contents) >= 3 {
		return contents[1], contents[2:]
	}

	return contents[1], nil

}

// Ensure that everything in a string array is lowercase
func lower(data []string) []string {
	var res []string

	for _, e := range data {
		res = append(res, strings.ToLower(e))
	}

	return res

}

// Link Twitch events to dedicated functions
func (g *Gambling) twitchOnEventSetup() {
	// Message handler, closure, because an access to *Gambling is needed
	g.Twitch.OnPrivateMessage(func(message twitch.PrivateMessage) {

		// the message does not contain the prefix, just return without doing nothing
		if !strings.HasPrefix(message.Message, g.Config.Prefix) {
			return
		}

		// Extract command and args from message
		cmd, args := extractCommand(message.Message)
		if cmd == "" {
			return
		}

		switch cmd {
		case "create":
			g.handleCreate(message.User, args)
			break
		case "close":
			g.handleClose(message.User)
			break
		case "roll":
			g.handleRoll(message.User, args)
			break
		case "vote":
			g.handleVote(message.User, args, g.Config.Verified)
			break
		case "delete":
			g.handleDelete(message.User)
			break
		case "reset":
			g.handleReset(message.User)
			break
		case "stats":
			g.handleStat(message.User, args)
			break
		default:
			fmt.Println("Warning: received an unsupported function")
		}

	})

	// On connect handler
	g.Twitch.OnConnect(func() {
		g.Twitch.Say(g.Config.Twitch.Channel, g.Config.Hello)
	})

}

// Join channel
func (g *Gambling) join() {
	g.Twitch.Join(g.Config.Twitch.Channel)
}

// Start is used to connect a gamble-bot instance to a twitch channel, and start the bot
func (g *Gambling) Start() error {
	return g.Twitch.Connect()
}

// say will be used to send informations to twitch channel
func (g *Gambling) say(message string) {
	g.Twitch.Say(g.Config.Twitch.Channel, message)
}

// say will be used to send informations to twitch channel
func (g *Gambling) sayAt(message string, users []string) {

	at := ""

	for _, u := range users {
		at = at + fmt.Sprintf("@%s ", u)
	}
	g.Twitch.Say(g.Config.Twitch.Channel, fmt.Sprintf("%s : %s", at, message))
}

// Handlers for all the things !
// create command handler
func (g *Gambling) handleCreate(user twitch.User, args []string) {

	// If user is not admin, return without doing nothing
	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	// Check if vote exists, by default, IsOpen will be false
	if g.CurrentVote.IsOpen {
		g.say(fmt.Sprintf("There is already a vote going, you should delete it first with '%s delete'.", g.Config.Prefix))
		return
	}

	if len(args) < 2 {
		g.say("You must pass the possibilities as arguments. (2 at least)")
		return
	}

	g.CurrentVote.IsOpen = true
	g.CurrentVote.Votes = make(map[string]string)
	g.CurrentVote.Possibilities = lower(args)

	g.say(fmt.Sprintf("There is a new vote! You can vote with '%s vote <vote> (choices are : %s)'", g.Config.Prefix, strings.Join(g.CurrentVote.Possibilities, " or ")))
}

// close vote handler
func (g *Gambling) handleClose(user twitch.User) {
	// If user is not admin, return without doing nothing
	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	g.CurrentVote.IsOpen = false
	g.say("The vote is now closed!")

}

// handle a vote
func (g *Gambling) handleVote(user twitch.User, args []string, verified bool) {

	// create a date value
	date := time.Now()

	// Ensure the vote is open
	if !g.CurrentVote.IsOpen {
		fmt.Println("Don't try to vote while it's close...")
		return
	}

	// Ensure there is args
	if args == nil || len(args) < 1 {
		fmt.Println("Expected one argument to vote...")
		return
	}

	// Ensure vote is lowercase
	vote := strings.ToLower(args[0])

	// Check if vote is valid
	if g.isVoteValid(vote) {
		// If it is add it
		g.CurrentVote.Votes[user.Name] = vote
		// ensure bot is verified to send whispers
		if verified {
			g.Twitch.Whisper(user.Name,
				fmt.Sprintf("For your information, I correctly handled your vote for %s (sent on %d-%02d-%02d)",
					vote, date.Year(), date.Month(), date.Day()))
		}
	} else {
		fmt.Println(user.Name + ": Invalid vote: got '" + args[0] + "'")
		// ensure bot is verified to send whispers
		if verified {
			g.Twitch.Whisper(user.Name, fmt.Sprintf("Option %s is invalid", vote))
		}
	}

}

// handle a vote delete
func (g *Gambling) handleDelete(user twitch.User) {

	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	g.CurrentVote = new(Vote)

}

// handle a vote reset
func (g *Gambling) handleReset(user twitch.User) {

	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	g.CurrentVote.Votes = make(map[string]string)
}

// handle a call to stats generation (public or private)
func (g *Gambling) handleStat(user twitch.User, args []string) {

	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	// create stats and store it into a string
	stats := createStat(g.CurrentVote)
	// Check args, if private mode is requested
	if len(args) >= 1 {
		if strings.ToLower(args[0]) == "private" {
			// Do it in private
			err := statsToFile(stats, g.Config.Stats.Dir)
			if err != nil {
				fmt.Println(err)
				g.say("Error generating stats in private mode")
			}
			g.say("Stats generated in private mode")
			return
		}
	} else {
		// If not in private mode, go public, and paste to pastebin
		link, err := statsToPastebin(g.Config.Pastebin.Key, stats)
		if err != nil {
			fmt.Println(err)
			g.say("Error generating stats in public mode")
		}
		g.say(fmt.Sprintf("Stats generated in public mode, check it at : %s", link))
	}

}

// handle roll and select winner
func (g *Gambling) handleRoll(user twitch.User, args []string) {

	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	if len(g.CurrentVote.Votes) == 0 {
		g.sayAt("There is no vote", g.Config.Admins)
		return
	}

	if g.CurrentVote.IsOpen {
		g.sayAt("The vote isn't closed", g.Config.Admins)
		return
	}

	if len(args) < 1 {
		g.sayAt("You must specify the winner", g.Config.Admins)
		return
	}

	if !g.isVoteValid(args[0]) {
		g.sayAt(fmt.Sprintf("This is not a correct option: %s", args[0]), g.Config.Admins)
	}

	winners := []string{}
	for user, vote := range g.CurrentVote.Votes {
		if vote == strings.ToLower(args[0]) {
			winners = append(winners, user)
		}
	}

	if len(winners) == 0 {
		g.sayAt("No one vote for this one...", g.Config.Admins)
		return
	}

	g.sayAt(fmt.Sprintf("The winner is... %s", winners[rand.Intn(len(winners))]), g.Config.Admins)

}

// Is a vote valid ?
func (g *Gambling) isVoteValid(vote string) bool {

	// Loop over possibilities
	for _, p := range g.CurrentVote.Possibilities {
		// ensure lowercase on everything
		// check equality
		if p == strings.ToLower(vote) {
			return true
		}
	}

	return false

}
