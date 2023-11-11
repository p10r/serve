package discord

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewFakeServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(400)
			return
		}
		w.WriteHeader(204)
	}))
}
