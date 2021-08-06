package reminders

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type cmd struct {
	action interface{}
	help   string
}

var commands = make(map[string]cmd)

func help() string {
	helpString := "Flags:\n"
	for f, c := range commands {
		helpString = helpString + fmt.Sprintf("-%v  %v\n", f, c.help)
	}
	return helpString
}

func ParseCommand(user discordgo.User, channel, command string) Reminder {
	// Declare variables
	var tags = []string{}
	var end time.Time
	var message string
	var flags = make(map[string]string)

	// split out the trailing message from the flags.
	// message must be wrapped in quotes
	cmd := strings.Split(command, "\"")
	message = cmd[1]

	// split the flags apart. Each should lead with a -
	flgs := strings.Split(cmd[0], "-")
	// loop through flags and split them into the flag and it's args storing them as a map
	// should be separated by spaces
	for _, f := range flgs {
		f := strings.Split(f, " ")
		flags[f[0]] = f[1]
	}

}
