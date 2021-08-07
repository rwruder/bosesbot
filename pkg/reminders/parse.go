package reminders

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type cmd struct {
	action func(string, *Reminder) error
	help   string
}

func setTime(s string, r *Reminder) error {
	var err error
	dur, err := time.ParseDuration(s)
	if (r.EndTime == time.Time{}) {
		r.EndTime = time.Now()
	}
	r.EndTime = r.EndTime.Add(dur)
	return err
}

func setDate(s string, r *Reminder) error {
	fmt.Println(s)

	var err error
	var returnTime time.Time
	now := time.Now()
	str := strings.Split(s, ".")
	fmt.Println(str)
	dateString := strings.Split(str[0], "/")
	var month, day, year int
	var hours, minutes int

	fmt.Println(str)
	switch len(str) {
	case 1:
		fmt.Println("Case 1")
		hours, minutes = now.Hour(), now.Minute()
	case 2:
		timeString := strings.Split(str[1], ":")
		if len(timeString) != 2 {
			err = fmt.Errorf("Time format should be hh:mm")
			return err
		}
		hours, err = strconv.Atoi(timeString[0])
		minutes, err = strconv.Atoi(timeString[1])
	default:
		err = fmt.Errorf("Format for date is m/d/yyyy hh:mm. Time is optional, but month, day, and year are required")
		return err
	}

	switch len(dateString) {
	case 2:
		year = now.Year()
	case 3:
		year, err = strconv.Atoi(dateString[2])
	default:
		err = fmt.Errorf("Format for date is m/d/yyyy.h:m. Time is optional, but month, day, and year are required")
		return err
	}
	month, err = strconv.Atoi(dateString[0])
	day, err = strconv.Atoi(dateString[1])
	if err != nil {
		return err
	}

	location, err := time.LoadLocation("Local")
	returnTime = time.Date(year, time.Month(month), day, hours, minutes, 0, 0, location)

	r.EndTime = returnTime
	return err
}

func setMentions(s string, r *Reminder) error {
	mentions := s
	var err error
	r.Mentions = mentions
	return err
}

var commands = map[string]cmd{
	"t": {
		action: setTime,
		help:   "Set the time a number of hours, minutes or seconds from now: 4h30m15s",
	},
	"d": {
		action: setDate,
		help:   "Set a reminder based on a set date. m/d/yyyy h:m. If year or time is left off it will default to the current one",
	},
	"m": {
		action: setMentions,
		help:   "If this flag is set everyone mentioned in the message will also be mentioned in the reminder",
	},
}

func help() string {
	helpString := "Flags:\n -h  displays this message.\n"
	for f, c := range commands {
		helpString = helpString + fmt.Sprintf("-%v  %v\n", f, c.help)
	}
	return helpString
}

func ParseCommand(user discordgo.User, channel, command string) (Reminder, error) {
	// Declare variables
	var reminder Reminder
	var err error
	var message string
	var flags = make(map[string]string)

	// split out the trailing message from the flags.
	// message must be wrapped in quotes
	cmd := strings.Split(command, "\"")
	message = cmd[1]
	reminder.User = user
	reminder.Channel = channel
	reminder.Message = message
	reminder.Mentions = ""

	// split the flags apart. Each should lead with a -
	flgs := strings.Split(cmd[0], "-")
	// loop through flags and split them into the flag and it's args storing them as a map
	// should be separated by spaces
	for _, f := range flgs {
		f := strings.Split(f, " ")
		flags[f[0]] = f[1]
	}

	if _, ok := flags["h"]; ok {
		err = fmt.Errorf(help())
		return reminder, err
	}
	for k, v := range flags {
		if c, ok := commands[k]; ok {
			err = c.action(v, &reminder)
		}
	}

	if reminder.EndTime.Before(time.Now()) {
		err = fmt.Errorf("%v is in the past.", reminder.EndTime)
	}

	return reminder, err
}
