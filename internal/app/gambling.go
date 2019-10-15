package app

import (
	"fmt"

	twitch "github.com/gempir/go-twitch-irc/v2"
)

// Gambling is a meta structure containing all the stuff needed by a Gambling instance
type Gambling struct {
	// Config from yaml file
	Config Conf
	// Twitch client
	Twitch *twitch.Client
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

	return g

}

// Link Twitch events to dedicated functions
func (g *Gambling) twitchOnEventSetup() {
	// TODO: move parseMessage and parseMessage related functions, in this package
	g.Twitch.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)
	})
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
