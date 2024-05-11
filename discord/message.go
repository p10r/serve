package discord

import (
	"fmt"
	"github.com/p10r/serve/domain"
	"strconv"
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

func NewMessage(matches domain.UntrackedMatches, currentTime time.Time) Message {
	date := currentTime.Format("Monday, 2 January 2006")

	var fields []Fields
	for league, matchesForCountry := range sortByLeague(matches) {
		fields = append(fields, Fields{
			Name:   flag(league) + matchesForCountry[0].FlashscoreName,
			Value:  text(matchesForCountry),
			Inline: false,
		})
	}

	return Message{fmt.Sprintf("Games for %s", date), []Embeds{{fields}}}
}

func sortByLeague(matches domain.UntrackedMatches) map[string]domain.UntrackedMatches {
	countries := make(map[string]domain.UntrackedMatches)
	for _, match := range matches {
		countries[match.FlashscoreName] = append(countries[match.FlashscoreName], match)
	}
	return countries
}

func text(matches domain.UntrackedMatches) string {
	var texts []string
	for _, e := range matches {
		texts = append(texts, fmt.Sprintf("**%v - %v**\t %v", e.HomeName, e.AwayName, hour(e.StartTime)))
	}
	return strings.Join(texts, "\n")
}

func flag(leagueName string) string {
	if strings.Contains(leagueName, "Poland") {
		return "🇵🇱"
	}
	if strings.Contains(leagueName, "Italy") {
		return "🇮🇹"
	}
	if strings.Contains(leagueName, "France") {
		return "🇫🇷"
	}
	if strings.Contains(leagueName, "Germany") {
		return "🇩🇪"
	}
	if strings.Contains(leagueName, "Russia") {
		return "🇷🇺"
	}
	if strings.Contains(leagueName, "Turkey") {
		return "🇹🇷"
	}
	if strings.Contains(leagueName, "Europe") {
		return "🇪🇺"
	}
	if strings.Contains(leagueName, "USA") {
		return "🇺🇸"
	}
	return ""
}

// See https://hammertime.cyou/ for more info
func hour(unixTs int64) string {
	return fmt.Sprintf("<t:%s:t>", strconv.FormatInt(unixTs, 10))
}
