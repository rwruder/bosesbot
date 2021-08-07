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

// Set the time based on hours, minutes, and seconds from now
// format 1h40m20s
func setTime(s string, r *Reminder) error {
	var err error
	dur, err := time.ParseDuration(s)
	if (r.EndTime == time.Time{}) {
		r.EndTime = time.Now()
	}
	r.EndTime = r.EndTime.Add(dur)
	return err
}

// set the reminder time based on a date
// string should be m/d/yyyy with an optional time h:m
// date and time should be separated by a .
func setDate(s string, r *Reminder) error {
	var err error
	var returnTime time.Time
	now := time.Now()
	str := strings.Split(s, ".")
	dateString := strings.Split(str[0], "/")
	var month, day, year int
	var hours, minutes int

	// if time was left off set time to the current time
	// if time was included set it to that time
	// if splitting the string results in something other than one or two strings
	// then an error should be raised.
	switch len(str) {
	case 1:
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

	// Check if year is included. If not use the current year
	// if splitting the date results in something other than 2 or 3 parts
	// then raise an error
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

	// Set location to local. Should change this to be user specifc in the future
	location, err := time.LoadLocation("Local")
	returnTime = time.Date(year, time.Month(month), day, hours, minutes, 0, 0, location)

	r.EndTime = returnTime
	return err
}

// sets the reminder's mentions to be s
func setMentions(s string, r *Reminder) error {
	mentions := s
	var err error
	r.Mentions = mentions
	return err
}

// Map of commands
// add commands here with the key as the flag
var commands = map[string]cmd{
	"t": {
		action: setTime,
		help:   "Set the time a number of hours, minutes or seconds from now: 4h30m15s",
	},
	"d": {
		action: setDate,
		help:   "Set a reminder based on a set date. m/d/yyyy.h:m. If year or time is left off it will default to the current one. Date and time should be separated by a period",
	},
	"m": {
		action: setMentions,
		help:   "should be a list of mentions. The mentions should not be separated by anything",
	},
}

// displays help for the various flags
func help() string {
	helpString := "Format: !r [dtm] \"(message)\"\n\nFlags:\n -h  displays this message.\n"
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
	// Split leaves an extra empty trailing string
	cmd := strings.Split(command, "\"")
	if len(cmd) >= 2 {
		message = cmd[1]
	} else {
		message = ""
	}
	// fill in the reminder so it's ready to be modified by the commands
	reminder.User = user
	reminder.Channel = channel
	reminder.Message = message
	reminder.Mentions = ""

	// split the flags apart. Each should lead with a -
	// remove the extra string before the -
	flgs := strings.Split(cmd[0], "-")[1:]
	// loop through flags and split them into the flag and it's args storing them as a map
	// should be separated by spaces, and there should be only
	for _, f := range flgs {
		// Separate flags by spaces and remove extra trailing string
		f := strings.Split(f, " ")[:len(f)]
		switch len(f) {
		case 1:
			flags[f[0]] = ""
		case 2:
			flags[f[0]] = f[1]
		default:
			err = fmt.Errorf("You entered the arguments to one of the flags incorrectly")
			return reminder, err
		}
	}

	// if -h is used return an error containing the help string
	// to stop the reminder from being used
	if _, ok := flags["h"]; ok {
		err = fmt.Errorf(help())
		return reminder, err
	}

	// Move through flags and compare them to commands
	// if flag coresponds to a command execute that command
	// if it does not raise an error
	for k, v := range flags {
		if c, ok := commands[k]; ok {
			err = c.action(v, &reminder)
		} else {
			err = fmt.Errorf("One or more of your flags was not recognized")
		}
	}

	// if nothing else has raised an error, but the reminder
	// is set for some time in the past raise an error
	if err == nil && reminder.EndTime.Before(time.Now()) {
		err = fmt.Errorf("%v is in the past.", reminder.EndTime)
	}

	return reminder, err
}
