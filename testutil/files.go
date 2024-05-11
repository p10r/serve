package testutil

import (
	json2 "encoding/json"
	"github.com/p10r/serve/db"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/flashscore"
	"io"
	"os"
	"testing"
)

func FlashscoreResponse(tb testing.TB) []byte {
	content, err := os.ReadFile("../testdata/flashscore-res.json")
	if err != nil {
		tb.Fail()
	}

	return content
}

func UntrackedMatches(tb testing.TB) domain.UntrackedMatches {
	var res flashscore.Response
	err := json2.Unmarshal(FlashscoreResponse(tb), &res)
	if err != nil {
		tb.Fail()
	}

	return res.ToUntrackedMatches()
}

func ReadFile(t *testing.T, filePath string) []byte {
	t.Helper()

	file, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	return data
}

// MustOpenDB returns a new, open DB. Fatal on error.
func MustOpenDB(tb testing.TB) *db.DB {
	tb.Helper()

	// Write to an in-memory database by default.
	// If the -dump flag is set, generate a temp file for the database.
	dsn := ":memory:"

	instance := db.NewDB(dsn)
	if err := instance.Open(); err != nil {
		tb.Fatal(err)
	}
	return instance
}
