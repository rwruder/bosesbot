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

const (
	Eyes = "\U0001F440"
)

// Variables used for command line parameters
var (
	Token string

	R                chan *reminders.Reminder
	active_reminders reminders.Active
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	R = make(chan *reminders.Reminder)
	ar, err := reminders.Load("save.json")
	if err != nil {
		active_reminders = reminders.Active{ActiveReminders: make(map[string]*reminders.Reminder)}
	} else {
		active_reminders = ar
	}
	active_reminders.RemoveOld()
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
	dg.AddHandler(messageReactionAdd)
	dg.AddHandler(messageReactionRemove)
	go reminders.Listen(dg, R)
	defer active_reminders.Save("save.json")

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	dg.Identify.Intents = discordgo.IntentsGuildMessages + discordgo.IntentsGuildMessageReactions

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
	if len(m.Content) < 3 {
		return
	}
	if m.Content[:2] != "!r" {
		return
	}

	reminder, err := reminders.ParseCommand(*m.Author, m.Message, m.ChannelID, m.Content[2:])
	if err != nil {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("%v", err),
		)
	} else {
		message, err := s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("%v set a reminder for %v. Click the eyes if you also want to be mentioned in this reminder.", reminder.User.Username, reminder.EndTime),
		)
		if err != nil {
			s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("%v", err),
			)
		} else {
			err := s.MessageReactionAdd(m.ChannelID, message.ID, "\U0001F440")
			if err != nil {
				fmt.Println(err)
			}
			active_reminders.ActiveReminders[message.ID] = &reminder
			go reminder.Set(R)
		}

	}
}

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	if reminder, ok := active_reminders.ActiveReminders[r.MessageID]; ok {
		if r.Emoji.Name == Eyes {
			user, err := s.User(r.UserID)
			if err != nil {
				fmt.Println(err)
			}
			reminder.Mentions = append(reminder.Mentions, *user)
		}
	}
}

func messageReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == s.State.User.ID {
		return
	}

	if reminder, ok := active_reminders.ActiveReminders[r.MessageID]; ok {
		if r.Emoji.Name == Eyes {
			user, err := s.User(r.UserID)
			if err != nil {
				fmt.Println(err)
			}

			for n, u := range reminder.Mentions {
				if user.ID == u.ID {
					// Remove the user from mentions
					reminder.Mentions[n] = reminder.Mentions[len(reminder.Mentions)-1]
					reminder.Mentions[len(reminder.Mentions)-1] = discordgo.User{}
					reminder.Mentions = reminder.Mentions[:len(reminder.Mentions)-1]
				}
			}
		}
	}
}
