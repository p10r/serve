package discord

import (
	"fmt"
	"github.com/p10r/serve/flashscore"
	"strings"
	"time"
)

type Message struct {
	Content string   `json:"content"`
	Embeds  []Embeds `json:"embeds"`
}

type Embeds struct {
	Fields []Fields `json:"fields"`
}

type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func NewMessage(leagues flashscore.Leagues) Message {
	currentTime := time.Now()
	date := currentTime.Format("Monday, 2 January 2006")

	var fields []Fields
	for _, league := range leagues {
		fields = append(fields, Fields{
			Name:   flag(league.Name),
			Value:  text(league.Events),
			Inline: false,
		})
	}

	return Message{fmt.Sprintf("Games for %s", date), []Embeds{{fields}}}
}

func text(events flashscore.Events) string {
	var texts []string
	for _, e := range events {
		texts = append(texts, fmt.Sprintf("**%v - %v**\t %v", e.HomeName, e.AwayName, hour(e.StartTime)))
	}
	return strings.Join(texts, "\n")
}

func flag(leagueName string) string {
	if strings.Contains(leagueName, "Poland") {
		return "🇵🇱 " + leagueName
	}
	if strings.Contains(leagueName, "Italy") {
		return "🇮🇹 " + leagueName
	}
	if strings.Contains(leagueName, "France") {
		return "🇫🇷 " + leagueName
	}
	if strings.Contains(leagueName, "Germany") {
		return "🇩🇪 " + leagueName
	}
	if strings.Contains(leagueName, "Russia") {
		return "🇷🇺 " + leagueName
	}
	if strings.Contains(leagueName, "Turkey") {
		return "🇹🇷 " + leagueName
	}
	if strings.Contains(leagueName, "Europe") {
		return "🇪🇺 " + leagueName
	}
	return leagueName
}

func hour(unixTs int64) string {
	ts := time.Unix(unixTs, 0)
	ts.Format("15:14")

	locations := []struct {
		name     string
		location *time.Location
	}{
		{"BER", locationOf("Europe/Berlin")},
		{"NY", locationOf("America/New_York")},
		{"LA", locationOf("America/Los_Angeles")},
		{"HK", locationOf("Asia/Hong_Kong")},
	}

	formattedTimes := make([]string, len(locations))
	for i, loc := range locations {
		localTime := ts.In(loc.location)
		formattedTimes[i] = localTime.Format("15:04")
	}

	formattedString := fmt.Sprintf("(%s %s/%s %s/%s %s/%s %s)",
		formattedTimes[0], locations[0].name,
		formattedTimes[1], locations[1].name,
		formattedTimes[2], locations[2].name,
		formattedTimes[3], locations[3].name,
	)

	return fmt.Sprintf(formattedString)
}

func locationOf(locationName string) *time.Location {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		panic(err) // Handle error appropriately
	}
	return loc
}
