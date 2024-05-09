package flashscore_test

import (
	"bytes"
	"fmt"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
	"github.com/p10r/serve/flashscore"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFlashscore(t *testing.T) {
	t.Run("deserializes flashscore response", func(t *testing.T) {
		json := expect.ReadFile(t, "../fixtures/flashscore-response.json")

		response, err := flashscore.NewResponse(io.NopCloser(bytes.NewBufferString(string(json))))
		expect.NoErr(t, err)

		fmt.Printf("%v", response)
		expected := flashscore.Response{
			Leagues: []flashscore.League{
				{
					"Austria: AVL",
					flashscore.Events{
						{
							HomeName:         "Sokol Wien",
							AwayName:         "VBK Klagenfurt",
							StartTime:        1697898600,
							HomeScoreCurrent: "3",
							AwayScoreCurrent: "1",
							Stage:            "FINISHED",
						},
						{
							HomeName:         "Wolfurt",
							AwayName:         "hotVolleys",
							StartTime:        1697965200,
							HomeScoreCurrent: "",
							AwayScoreCurrent: "",
							Stage:            "SCHEDULED",
						},
					},
				},
			},
		}

		expect.True(t, reflect.DeepEqual(response, expected))
	})

	t.Run("returns error if json cannot be read", func(t *testing.T) {

	})

	t.Run("fetches response", func(t *testing.T) {
		apiKey := "1234"
		flashscoreServer := flashscore.NewFakeServer(t, apiKey)
		defer flashscoreServer.Close()

		client := flashscore.NewClient(flashscoreServer.URL, apiKey)

		_, err := client.GetUpcomingMatches()
		expect.NoErr(t, err)
	})

	t.Run("reports error", func(t *testing.T) {
		flashscoreServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
		}))
		defer flashscoreServer.Close()

		client := flashscore.NewClient(flashscoreServer.URL, "apiKey")

		_, err := client.GetUpcomingMatches()
		expect.Err(t, err)
	})

	t.Run("map flashscore response to flashscore matches", func(t *testing.T) {
		resp := flashscore.Response{
			Leagues: flashscore.Leagues{flashscore.League{
				Name: "Germany: 1. Bundesliga",
				Events: flashscore.Events{
					flashscore.Event{
						HomeName:         "Berlin",
						AwayName:         "Haching",
						StartTime:        1234,
						HomeScoreCurrent: "3",
						AwayScoreCurrent: "1",
						Stage:            "FINISHED",
					},
				},
			}}}

		expected := domain.UntrackedMatches{
			domain.UntrackedMatch{
				HomeName:       "Berlin",
				AwayName:       "Haching",
				StartTime:      1234,
				FlashscoreName: "Germany: 1. Bundesliga",
				Country:        "Germany",
				League:         "1. Bundesliga",
				Stage:          "FINISHED",
			},
		}

		expect.DeepEqual(t, resp.ToUntrackedMatches(), expected)
	})

}
