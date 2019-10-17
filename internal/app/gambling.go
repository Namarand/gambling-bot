package app

import (
	"fmt"
	"strings"

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
			// handleRoll(message.User, contents)
			break
		case "vote":
			g.handleVote(message.User, args)
			break
		case "delete":
			g.handleDelete(message.User)
			break
		case "reset":
			g.handleReset(message.User)
			break
		case "stats":
			// handleStat(message.User)
			break
		case "privatestats":
			// handlePrivateStat(message.User)
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

// say will be used to send informations to targeted users
// ENSURE BOT IS VERIFIED TO USE THIS FUNCTION
func (g *Gambling) whisper(message string, users []string) {
	for _, u := range users {
		g.Twitch.Whisper(u, message)
	}
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
func (g *Gambling) handleVote(user twitch.User, args []string) {

	// Ensure the vote is open
	if !g.CurrentVote.IsOpen {
		fmt.Println("Don't try to vote while it's close...")
		return
	}

	// Ensure there is args
	if args != nil && len(args) < 1 {
		fmt.Println("Expected one argument to vote...")
	}

	// Check if vote is valid
	if g.isVoteValid(args[0]) {
		// If it is add it
		g.CurrentVote.Votes[user.Name] = args[0]
	} else {
		fmt.Println(user.Name + ": Invalid vote: got '" + args[0] + "'")
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
