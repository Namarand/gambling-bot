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
			// handleClose(message.User)
			break
		case "roll":
			// handleRoll(message.User, contents)
			break
		case "vote":
			// handleVote(message.User, contents)
			break
		case "delete":
			// handleDelete(message.User)
			break
		case "reset":
			// handleReset(message.User)
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
	g.CurrentVote.Possibilities = args

	g.say(fmt.Sprintf("There is a new vote! You can vote with '%s vote <vote> (choices are : %s)", g.Config.Prefix, strings.Join(g.CurrentVote.Possibilities, " or ")))
}
