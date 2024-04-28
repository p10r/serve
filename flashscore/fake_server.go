package flashscore

import (
	"github.com/p10r/serve/expect"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

		json := expect.ReadFile(t, "../helpers/flashscore-response.json")
		w.Write([]byte(json))
	}))
}
