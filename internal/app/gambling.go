package app

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"context"

	"golang.org/x/time/rate"

	"github.com/apex/log"
	twitch "github.com/gempir/go-twitch-irc/v2"
)

// Size ack buffered chan
const bufferSize = 500

// Vote is a structure handling all voting params and status
type Vote struct {
	IsOpen        bool
	Possibilities []string
	Votes         map[string]string
	Acks          Acks
}

// Acks is used to store and send ack messages stored if rate limit is reached
type Acks struct {
	Buffer chan VoteAck
	WIP    bool
	Drop   chan bool
}

// NewAcks is used to init a Acks struct
func NewAcks() Acks {
	return Acks{
		Buffer: make(chan VoteAck, bufferSize),
		WIP:    false,
		Drop:   make(chan bool),
	}
}

// VoteAck is a struct containing vote acknolegment send later if rate limit is reached
type VoteAck struct {
	message  string
	Username string
}

// NewVoteAck is used to init a VoteAck struct
func NewVoteAck(message string, username string) VoteAck {
	return VoteAck{
		message:  message,
		Username: username,
	}
}

// Gambling is a meta structure containing all the stuff needed by a Gambling instance
type Gambling struct {
	// Config from yaml file
	Config Conf
	// Twitch client
	Twitch *twitch.Client
	// Vote
	CurrentVote *Vote
	// whisper rate limiter
	WhispRL *rate.Limiter
	// Warning rate limiter
	WarnRL *rate.Limiter
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

	// setup rate limiter for whispers
	// 19 times per second, burst set to 1
	g.WhispRL = rate.NewLimiter(19, 1)

	// setup rate limiter for warning messages
	// One every 15 seconds
	g.WarnRL = rate.NewLimiter(rate.Every(15*time.Second), 1)

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
			log.WithField("command", cmd).Warn("Unsupported command received")
			g.say("Sorry but this is not a supported command")
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
	log.Info("Bot instance is starting")
	return g.Twitch.Connect()
}

// say will be used to send informations to twitch channel
func (g *Gambling) say(message string) {
	g.Twitch.Say(g.Config.Twitch.Channel, message)
}

// choices function is used to return all possibilites in a vote as a string
func (g *Gambling) choices() string {
	return strings.Join(g.CurrentVote.Possibilities, " or ")
}

// whisper will be used to send whisper to some user with rate limit constraints
func (g *Gambling) whisper(user string, message string) error {

	// setup a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// ensure cancel
	defer cancel()
	// wait with context
	err := g.WhispRL.Wait(ctx)

	// check for errors
	if err != nil {
		// setup a context with a timeout
		// do not wait for a long time
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		// ensure cancel
		defer cancel()
		// wait with context
		err := g.WarnRL.Wait(ctx)
		if err != nil {
			// log as warning
			log.WithFields(log.Fields{
				"message": message,
				"user":    user,
			}).Warn("Rate limit reached, private message not sent")
			return errors.New("Can not send whisper nor warning message, rate limit reached")
		}

		g.say("Warning, whisper rate limit reached, you may receive your vote acknowledgement later")

		return errors.New("Can not send whisper message, rate limit reached, warning message send")
	}

	// if rate limit not reach, send message
	g.Twitch.Whisper(user, message)

	// logs message send as Info
	log.WithFields(log.Fields{
		"message": message,
		"user":    user,
	}).Info("Private message send")

	return err
}

// sayAt will be used to send messages to twitch channel with mentions to users passed as arguments
func (g *Gambling) sayAt(message string, users []string) {

	at := ""

	for _, u := range users {
		at = at + fmt.Sprintf("@%s ", u)
	}
	g.Twitch.Say(g.Config.Twitch.Channel, fmt.Sprintf("%s : %s", at, message))
}

// ackMessage is used to generated an ack message
func ackMessage(valid bool, vote string) string {
	// get current date
	date := time.Now()

	// tail of the message, used to specify sending date, since Twitch UI not clear about this
	tail := fmt.Sprintf("(sent on %d-%02d-%02d)", date.Year(), date.Month(), date.Day())

	// return a message for a valid vote
	if valid {
		return fmt.Sprintf("For your information, I correctly handled your vote for %s", vote) + " " + tail
	}

	// return a kind error message if vote if note valid
	return "Sorry but the vote command you send is not valid, you may have made a mistake, please retry" + " " + tail

}

// SendAcks is used to send Ack accumulted while rate limit is reached
func (g *Gambling) SendAcks() {

	// Mark Ack sending jobs as pending
	g.CurrentVote.Acks.WIP = true

	// loop until a drop is received
	for {
		select {
		// if a Drop inscrution is received, leave
		case <-g.CurrentVote.Acks.Drop:
			log.Info("Message drop instruction received")
			return
		// if an ack is received handle it
		case ack, ok := <-g.CurrentVote.Acks.Buffer:
			// if channel is empty juste leave
			if !ok {
				// Mark ack sending jobs as done
				log.Info("Acks WIP set to false")
				g.CurrentVote.Acks.WIP = false
				return
			}

			// Ensure a sleep of half a second between each messages, in order to be sure to not reach rate limit
			time.Sleep((200 * time.Millisecond))

			// try sending message from acks queue
			err := g.whisper(ack.Username, ack.message)
			if err != nil {
				// If rate limit is reached drop and log
				log.WithFields(log.Fields{
					"message": ack.message,
					"user":    ack.Username,
				}).Warn("Rate limit reached, while sending late ACKs message dropped")
				return
			}
		}
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
		g.say("You need to pass the choices as arguments (2 at least)")
		return
	}

	g.CurrentVote.IsOpen = true
	g.CurrentVote.Votes = make(map[string]string)
	g.CurrentVote.Possibilities = lower(args)
	// Ensure a 500 items long ACK queue
	g.CurrentVote.Acks = NewAcks()

	g.say(fmt.Sprintf("There is a new vote! You can vote with '%s vote <vote> (choices are : %s)'", g.Config.Prefix, g.choices()))

	log.WithFields(log.Fields{
		"choices":        g.CurrentVote.Possibilities,
		"requested by":   user.DisplayName,
		"acks queue len": bufferSize,
	}).Info("Vote created")
}

// close vote handler
func (g *Gambling) handleClose(user twitch.User) {
	// If user is not admin, return without doing nothing
	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	// Check if vote is open
	if !g.CurrentVote.IsOpen {
		g.say("Do not try to close an alreay closed vote !")
		return
	}

	g.CurrentVote.IsOpen = false

	// close channel
	close(g.CurrentVote.Acks.Buffer)

	// unpile, if verified
	if g.Config.Verified {
		go g.SendAcks()
	}

	st := NewStatistics(g.CurrentVote)

	var parts []string

	for k, v := range st.Transformed {
		parts = append(parts, fmt.Sprintf("%s : %d (%.2f%%)", k, len(v), (float64(len(v))/float64(st.Total))*100))
	}

	log.Info("Vote closed")

	g.say("Vote is now closed, time for statistics ! " + fmt.Sprintf("Participants : %d", st.Total) + " | " + strings.Join(parts, ", "))

}

// handle a vote
func (g *Gambling) handleVote(user twitch.User, args []string, verified bool) {

	// Ensure the vote is open
	if !g.CurrentVote.IsOpen {
		log.Warn("Vote triggered while close")
		return
	}

	// Ensure there is args
	if args == nil || len(args) < 1 {
		if verified {
			message := ackMessage(false, "")
			err := g.whisper(user.Name, message)
			// if an error occur, rate limit is reached try later and add it to ack queue
			if err != nil {
				log.WithFields(log.Fields{
					"message": message,
					"user":    user,
					"valid":   false,
				}).Info("Message added to ACKs queue")
				g.CurrentVote.Acks.Buffer <- NewVoteAck(message, user.Name)
			}
		}
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
			message := ackMessage(true, vote)
			err := g.whisper(user.Name, message)
			// if an error occur, rate limit is reached try later and add it to ack queue
			if err != nil {
				log.WithFields(log.Fields{
					"message": message,
					"user":    user,
					"valid":   true,
				}).Info("Message added to ACKs queue")
				g.CurrentVote.Acks.Buffer <- NewVoteAck(message, user.Name)
			}
		}
	} else {
		// ensure bot is verified to send whispers
		if verified {
			message := ackMessage(false, "")
			err := g.whisper(user.Name, message)
			// if an error occur, rate limit is reached try later and add it to ack queue
			if err != nil {
				log.WithFields(log.Fields{
					"message": message,
					"user":    user,
					"valid":   false,
				}).Info("Message added to ACKs queue")
				g.CurrentVote.Acks.Buffer <- NewVoteAck(message, user.Name)
			}
		}
	}

}

// handle a vote delete
func (g *Gambling) handleDelete(user twitch.User) {

	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	// if there is a working job sending acks
	if g.CurrentVote.Acks.WIP {
		// Drop all jobs
		g.CurrentVote.Acks.Drop <- true
	}

	// Create a new empty vote
	g.CurrentVote = new(Vote)

	log.Info("Vote deleted")

	g.say("Vote deleted !")

}

// handle a vote reset
func (g *Gambling) handleReset(user twitch.User) {

	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	// if there is a working job sending acks
	if g.CurrentVote.Acks.WIP {
		// Drop all jobs
		g.CurrentVote.Acks.Drop <- true
	}

	g.CurrentVote.Acks.Buffer = make(chan VoteAck, bufferSize)

	g.CurrentVote.Votes = make(map[string]string)

	log.Info("Vote reset")
}

// handle a call to stats generation (public or private)
func (g *Gambling) handleStat(user twitch.User, args []string) {

	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	// create stats and store it into a string
	stats := createStat(g.CurrentVote)
	err := statsToFile(stats, g.Config.Stats.Dir)
	if err != nil {
		log.Error(err.Error())
		g.say("Error generating statistics")
	}
	g.say("Statistics generated in private mode")

}

// handle roll and select winner
func (g *Gambling) handleRoll(user twitch.User, args []string) {

	if !checkPermission(user.Name, g.Config.Admins) {
		return
	}

	if len(g.CurrentVote.Votes) == 0 {
		g.sayAt("You can not roll since there is no vote", g.Config.Admins)
		return
	}

	if g.CurrentVote.IsOpen {
		g.sayAt(fmt.Sprintf("Hey ! The vote isn't closed ! Close it using command : '%s close'", g.Config.Prefix), g.Config.Admins)
		return
	}

	if len(args) < 1 {
		g.sayAt(fmt.Sprintf("You must specify the winner (choices are : %s)", g.choices()), g.Config.Admins)
		return
	}

	if !g.isVoteValid(args[0]) {
		g.sayAt(fmt.Sprintf("%s is not a correct roll option (choices are : %s)", args[0], g.choices()), g.Config.Admins)
	}

	winners := []string{}
	for user, vote := range g.CurrentVote.Votes {
		if vote == strings.ToLower(args[0]) {
			winners = append(winners, user)
		}
	}

	if len(winners) == 0 {
		log.WithField("selected", args[0]).Warn("No one vote for this one")
		g.sayAt("Sorry but, no one vote for this one...", g.Config.Admins)
		return
	}

	winner := winners[rand.Intn(len(winners))]

	g.sayAt(fmt.Sprintf("And... The winner is... %s", winner), g.Config.Admins)

	log.WithField("user", winner).Warn("Randomly selected user")

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
