package flashscore_test

import (
	"bytes"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
	"github.com/p10r/serve/flashscore"
	"github.com/p10r/serve/testutil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO move to testutil
func NewFakeServer(t *testing.T, apiKey string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(400)
			return
		}

		reqApiKey := r.Header.Get("X-RapidAPI-Key")
		if reqApiKey != apiKey {
			w.WriteHeader(400)
			return
		}

		json := testutil.FlashscoreResponse(t)
		w.Write([]byte(json))
	}))
}

func TestFlashscore(t *testing.T) {
	t.Run("deserializes flashscore response", func(t *testing.T) {
		json := testutil.FlashscoreResponse(t)

		response, err := flashscore.NewResponse(io.NopCloser(bytes.NewBufferString(string(json))))

		expect.NoErr(t, err)
		expect.Equal(t, len(response.Leagues), 5)
	})

	t.Run("fetches response", func(t *testing.T) {
		apiKey := "1234"
		flashscoreServer := NewFakeServer(t, apiKey)
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
