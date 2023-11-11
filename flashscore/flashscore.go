package flashscore

import (
	json2 "encoding/json"
	"io"
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
