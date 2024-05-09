package flashscore

import (
	json2 "encoding/json"
	"github.com/p10r/serve/domain"
	"io"
	"strings"
)

type Response struct {
	Leagues Leagues `json:"DATA"`
}

func NewResponse(input io.ReadCloser) (Response, error) {
	var res Response

	err := json2.NewDecoder(input).Decode(&res)
	if err != nil {
		return Response{}, err
	}

	return res, nil
}

func (r Response) ToUntrackedMatches() domain.UntrackedMatches {
	matches := domain.UntrackedMatches{}

	for _, league := range r.Leagues {
		leagueInfo := strings.Split(league.Name, ": ")
		country := leagueInfo[0]
		leagueName := leagueInfo[1]

		for _, event := range league.Events {
			match := domain.UntrackedMatch{
				HomeName:       event.HomeName,
				AwayName:       event.AwayName,
				StartTime:      event.StartTime,
				FlashscoreName: league.Name,
				Country:        country,
				League:         leagueName,
				Stage:          event.Stage,
			}

			matches = append(matches, match)
		}
	}

	return matches
}

type Leagues []League

type League struct {
	Name   string `json:"NAME"`
	Events Events `json:"EVENTS"`
}

type Events []Event

type Event struct {
	HomeName         string `json:"HOME_NAME"`
	AwayName         string `json:"AWAY_NAME"`
	StartTime        int64  `json:"START_TIME"`
	HomeScoreCurrent string `json:"HOME_SCORE_CURRENT"`
	AwayScoreCurrent string `json:"AWAY_SCORE_CURRENT"`
	Stage            string `json:"STAGE"`
}
