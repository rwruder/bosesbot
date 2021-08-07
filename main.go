package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rwruder/bosesbot/pkg/reminders"
)

// Variables used for command line parameters
var (
	Token string
	R     chan *reminders.Reminder
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	R = make(chan *reminders.Reminder)
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(reminderCreate)
	go reminders.Listen(dg, R)

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func reminderCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content[:2] != "!r" {
		return
	}

	reminder, err := reminders.ParseCommand(*m.Author, m.ChannelID, m.Content[2:])
	if err != nil {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("%v", err),
		)
	} else {
		go reminder.Set(R)
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("%v set a reminder for %v.", reminder.User.Username, reminder.EndTime),
		)
	}
}
