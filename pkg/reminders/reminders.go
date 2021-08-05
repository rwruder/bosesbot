package reminders

import "time"

type Reminder struct {
	User    string    `json:"user"`
	Tags    []string  `json:"tags"`
	Channel string    `json:"channel"`
	EndTime time.Time `json:"end"`
	Message string    `json:"message"`
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
