package reminders

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Reminder struct {
	User    discordgo.User `json:"user"`
	Tags    []string       `json:"tags"`
	Channel string         `json:"channel"`
	EndTime time.Time      `json:"end"`
	Message string         `json:"message"`
}

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
