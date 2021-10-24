package reminders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Active struct {
	ActiveReminders map[string]*Reminder `json:"activereminders"`
}

type Reminder struct {
	User           discordgo.User     `json:"user"`
	Channel        string             `json:"channel"`
	Mentions       []discordgo.User   `json:"mentions"`
	DiscordMessage *discordgo.Message `json:"discordmessage"`
	EndTime        time.Time          `json:"end"`
	Message        string             `json:"message"`
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
		mention := done.User.Mention()
		for _, m := range done.Mentions {
			mention += m.Mention()
		}
		message := done.Message
		remind := fmt.Sprintf("%v %v", mention, message)
		s.ChannelMessageSend(done.Channel, remind)
	}
}

func (a *Active) Save(file string) error {
	s, err := json.Marshal(a.ActiveReminders)
	if err != nil {
		return err
	}
	ioutil.WriteFile(file, s, 0644)
	return err
}

func (a *Active) RemoveOld() {
	for key, act := range a.ActiveReminders {
		if act.EndTime.After(time.Now()) {
			delete(a.ActiveReminders, key)
		}
	}
}

func Load(file string) (Active, error) {
	var active Active
	a, err := ioutil.ReadFile(file)
	if err != nil {
		return active, err
	}
	json.Unmarshal(a, &active)
	return active, err
}
