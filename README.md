This is my discord bot. All parts of it are currently a work in progress

######packages

reminders - This should allow the bot to be used to set reminder messages. At the designated time it will @ the user and other selected users or roles and display a message

Usage: 
Reminder commands all start with !r and contains a number of flags followed by a message contained in quotes. The bot will display a message confirming the reminder and then will wait the duration and then send the message specified in the command. It will always tag the user who made the reminder, but can also be set to mention other users with the -m command.

There are two methods to set the date and time of the reminder -t and -d. -t sets the reminder for a number of hours, minutes in the future, and -d will set it for a specific date and time.

for more information on the various flags and specific formats use !r -h for a help message from the bot