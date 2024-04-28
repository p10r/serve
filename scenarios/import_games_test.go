package scenarios_test

import (
	"context"
	"github.com/p10r/serve/db"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/helpers"
	"testing"
)

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

func TestImportMatches(t *testing.T) {
	ctx := context.TODO()

	t.Run("imports today's matches to db", func(t *testing.T) {
		//MatchImporter{}.importMatches()

		untrackedMatch := domain.UntrackedMatch{
			HomeName:  "Berlin",
			AwayName:  "Düren",
			StartTime: 123,
			Country:   "Germany",
			League:    "Bundesliga Playoffs",
		}

		matchStore := db.NewMatchStore(MustOpenDB(t))

		_, err := matchStore.Add(ctx, untrackedMatch)
		helpers.NoErr(t, err)

		matches, err := matchStore.All(ctx)
		helpers.NoErr(t, err) //TODO rename helpers to expect

		expected := domain.Match{
			ID:        1,
			HomeName:  "Berlin",
			AwayName:  "Düren",
			StartTime: 123,
			Country:   "Germany",
			League:    "Bundesliga Playoffs",
		}
		helpers.Equal(t, matches[0], expected)
	})

	//test what happens if two matches with the same timestamp are in db
}
