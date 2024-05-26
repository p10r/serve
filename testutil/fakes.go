package testutil

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DiscordServer struct {
	*httptest.Server
	Requests *[][]byte
}

func NewDiscordServer(t *testing.T) *DiscordServer {
	t.Helper()

	var reqs [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			log.Println("DiscordServer: Received invalid request")
			w.WriteHeader(400)
			return
		}

		log.Println("DiscordServer: Received request")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Cannot read body")
			return
		}
		defer r.Body.Close()
		reqs = append(reqs, body)

		w.WriteHeader(204)
	}))

	return &DiscordServer{server, &reqs}
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

		body, err := json.Marshal(FlashscoreRes(t))
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
