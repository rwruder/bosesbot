package reminders_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rwruder/bosesbot/pkg/reminders"

	"github.com/bwmarrin/discordgo"
)

func acceptableRange(t, target time.Time, pad time.Duration) bool {
	return t.Before(target.Add(pad)) && t.After(target.Add(-pad))
}

func TestParseCommandBasicTimer(t *testing.T) {
	input := " -t 3m"
	target := time.Now().Add(time.Minute * 3)
	t_remind, err := reminders.ParseCommand(discordgo.User{}, "", input)
	if err != nil {
		t.Errorf("Basic Timer threw an error: %v", err)
	}
	if !acceptableRange(t_remind.EndTime, target, time.Second) {
		diff := t_remind.EndTime.Sub(target).String()
		t.Errorf("Time was %v off from the expected value", diff)
	}

}

func TestParseCommandBasicDate(t *testing.T) {
	now := time.Now()
	year, month, day, hour, minute := now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute()
	year += 1
	month += 1
	day += 1
	hour += 1
	minute += 1
	input := fmt.Sprintf("-d %v/%v/%v.%v:%v", month, day, year, hour, minute)
	target := now.AddDate(1, 1, 1).Add(time.Hour * 1).Add(time.Minute * 1)
	t_remind, err := reminders.ParseCommand(discordgo.User{}, "", input)
	if err != nil {
		t.Errorf("Basic Date threw an error: %v", err)
	}
	if !acceptableRange(t_remind.EndTime, target, time.Minute) {
		t.Errorf("Expected %v for time, but instead got %v", target, t_remind.EndTime)
	}
}

func TestParseCommandTimerMessage(t *testing.T) {
	input := " -t 3m \"Test Message\""
	target_t := time.Now().Add(time.Minute * 3)
	target_m := "Test Message"
	t_remind, err := reminders.ParseCommand(discordgo.User{}, "", input)
	if err != nil {
		t.Errorf("Basic Timer threw an error: %v", err)
	}
	if !acceptableRange(t_remind.EndTime, target_t, time.Second) {
		diff := t_remind.EndTime.Sub(target_t).String()
		t.Errorf("Time was %v off from the expected value", diff)
	}
	if t_remind.Message != target_m {
		t.Errorf("Expected message to be %v, but was %v", target_m, t_remind.Message)
	}
}

func TestParseCommandDateMessage(t *testing.T) {
	now := time.Now()
	year, month, day, hour, minute := now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute()
	year += 1
	month += 1
	day += 1
	hour += 1
	minute += 1
	input := fmt.Sprintf("-d %v/%v/%v.%v:%v \"Test Message\"", month, day, year, hour, minute)
	target_t := now.AddDate(1, 1, 1).Add(time.Hour * 1).Add(time.Minute * 1)
	target_m := "Test Message"
	t_remind, err := reminders.ParseCommand(discordgo.User{}, "", input)
	if err != nil {
		t.Errorf("Basic Date threw an error: %v", err)
	}
	if !acceptableRange(t_remind.EndTime, target_t, time.Minute) {
		t.Errorf("Expected %v for time, but instead got %v", target_t, t_remind.EndTime)
	}
	if t_remind.Message != target_m {
		t.Errorf("Expected message to be %v, but was %v", target_m, t_remind.Message)
	}
}

func TestParseCommandTimerMentions(t *testing.T) {
	input := "-t 3m -m @bosesbjorn"
	target_t := time.Now().Add(time.Minute * 3)
	target_m := "@bosesbjorn"
	t_remind, err := reminders.ParseCommand(discordgo.User{}, "", input)
	if err != nil {
		t.Errorf("Basic Timer threw an error: %v", err)
	}
	if !acceptableRange(t_remind.EndTime, target_t, time.Second) {
		diff := t_remind.EndTime.Sub(target_t).String()
		t.Errorf("Time was %v off from the expected value", diff)
	}
	if t_remind.Mentions != target_m {
		t.Errorf("Expected mentions to be %v but actually got %v", target_m, t_remind.Mentions)
	}

}

func TestParseCommandDateMentions(t *testing.T) {
	now := time.Now()
	year, month, day, hour, minute := now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute()
	year += 1
	month += 1
	day += 1
	hour += 1
	minute += 1
	input := fmt.Sprintf("-d %v/%v/%v.%v:%v -m @bosesbjorn", month, day, year, hour, minute)
	target_t := now.AddDate(1, 1, 1).Add(time.Hour * 1).Add(time.Minute * 1)
	target_m := "@bosesbjorn"
	t_remind, err := reminders.ParseCommand(discordgo.User{}, "", input)
	if err != nil {
		t.Errorf("Basic Date threw an error: %v", err)
	}
	if !!acceptableRange(t_remind.EndTime, target_t, time.Minute) {
		t.Errorf("Expected %v for time, but instead got %v", target_t, t_remind.EndTime)
	}
	if t_remind.Mentions != target_m {
		t.Errorf("Expected mentions to be %v, but was %v", target_m, t_remind.Mentions)
	}
}

func TestParseCommandTimerMentionsMessage(t *testing.T) {
	input := "-t 3m -m @bosesbjorn \"Test Message\""
	target_t := time.Now().Add(time.Minute * 3)
	target_men := "@bosesbjorn"
	target_mes := "Test Message"
	t_remind, err := reminders.ParseCommand(discordgo.User{}, "", input)
	if err != nil {
		t.Errorf("Basic Timer threw an error: %v", err)
	}
	if !acceptableRange(t_remind.EndTime, target_t, time.Second) {
		diff := t_remind.EndTime.Sub(target_t).String()
		t.Errorf("Time was %v off from the expected value", diff)
	}
	if t_remind.Mentions != target_men {
		t.Errorf("Expected mentions to be %v but actually got %v", target_men, t_remind.Mentions)
	}
	if t_remind.Message != target_mes {
		t.Errorf("Expected message to be %v, but was %v", target_mes, t_remind.Message)
	}
}

func TestParseCommandDateMentionsMessage(t *testing.T) {
	now := time.Now()
	year, month, day, hour, minute := now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute()
	year += 1
	month += 1
	day += 1
	hour += 1
	minute += 1
	input := fmt.Sprintf("-d %v/%v/%v.%v:%v -m @bosesbjorn \"Test Message\"", month, day, year, hour, minute)
	target_t := now.AddDate(1, 1, 1).Add(time.Hour * 1).Add(time.Minute * 1)
	target_men := "@bosesbjorn"
	target_mes := "Test Message"
	t_remind, err := reminders.ParseCommand(discordgo.User{}, "", input)
	if err != nil {
		t.Errorf("Basic Date threw an error: %v", err)
	}
	if !acceptableRange(t_remind.EndTime, target_t, time.Minute) {
		t.Errorf("Expected %v for time, but instead got %v", target_t, t_remind.EndTime)
	}
	if t_remind.Mentions != target_men {
		t.Errorf("Expected mentions to be %v, but was %v", target_men, t_remind.Mentions)
	}
	if t_remind.Message != target_mes {
		t.Errorf("Expected message to be %v, but was %v", target_mes, t_remind.Message)
	}
}

func TestParseCommandErrors(t *testing.T) {

}
