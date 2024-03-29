package flashscore_test

import (
	"bytes"
	"fmt"
	"github.com/p10r/serve/flashscore"
	"github.com/p10r/serve/helpers"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFlashscore(t *testing.T) {
	t.Run("deserializes flashscore response", func(t *testing.T) {
		json := helpers.ReadFile(t, "../helpers/flashscore-response.json")

		response, err := flashscore.NewResponse(io.NopCloser(bytes.NewBufferString(string(json))))
		helpers.NoErr(t, err)

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

		helpers.True(t, reflect.DeepEqual(response, expected))
	})

	t.Run("returns error if json cannot be read", func(t *testing.T) {

	})

	t.Run("fetches response", func(t *testing.T) {
		apiKey := "1234"
		flashscoreServer := flashscore.NewFakeServer(t, apiKey)
		defer flashscoreServer.Close()

		client := flashscore.NewClient(flashscoreServer.URL, apiKey)

		_, err := client.GetSchedule()
		helpers.NoErr(t, err)
	})

	t.Run("reports error", func(t *testing.T) {
		flashscoreServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
		}))
		defer flashscoreServer.Close()

		client := flashscore.NewClient(flashscoreServer.URL, "apiKey")

		_, err := client.GetSchedule()
		helpers.Err(t, err)
	})
}
