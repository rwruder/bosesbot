package reminders

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Reminder struct {
	User     discordgo.User `json:"user"`
	Mentions string         `json:"mentions"`
	Channel  string         `json:"channel"`
	EndTime  time.Time      `json:"end"`
	Message  string         `json:"message"`
}

// Starts the reminder
func (r *Reminder) Set(end chan *Reminder) {
	// Takes a channel and sends the reminder to that channel when the timer is up
	for {
		timer := time.NewTimer(time.Until(r.EndTime))
		done := <-timer.C
		if time.Now().After(done) {
			break
		}
	}
	end <- r
}

// Starts listening for reminders on channel r
func Listen(s *discordgo.Session, r chan *Reminder) {
	for {
		done := <-r
		mention := done.User.Mention() + done.Mentions
		message := done.Message
		remind := fmt.Sprintf("%v %v", mention, message)
		s.ChannelMessageSend(done.Channel, remind)
	}
}
