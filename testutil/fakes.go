package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewDiscordServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(400)
			return
		}
		w.WriteHeader(204)
	}))
}

func NewFlashscoreServer(t *testing.T, apiKey string) *httptest.Server {
	t.Helper()

	//https://flashscore.p.rapidapi.com/v1/events/list?locale=en_GB&timezone=-4&sport_id=12&indent_days=0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(400)
			t.Fatalf("Flashscore Server: Invalid req method")
			return
		}

		if r.URL.Path != "/v1/events/list" {
			w.WriteHeader(400)
			t.Fatalf("Flashscore Server: Invalid URL path")
			return
		}

		apiKeyHeader := r.Header.Get("X-RapidAPI-Key")
		if apiKeyHeader != apiKey {
			w.WriteHeader(400)
			t.Fatalf("Flashscore Server: X-RapidAPI-Key does not match. \n\t\tGot: %v \n\t\tWant: %v", apiKeyHeader, apiKey)
			return
		}

		res := FlashscoreResponse(t)
		body, err := json.Marshal(res)
		if err != nil {
			t.Fatal("could not marshall JSON")
		}
		// TODO: check for X-RapidAPI-Host and X-RapidAPI-Key
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(body)
		if err != nil {
			t.Fatalf("could not set response: %v", err)
		}
	}))
}
