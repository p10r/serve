package scenarios_test

import (
	"context"
	"github.com/p10r/serve/db"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
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
		untrackedMatch := domain.UntrackedMatch{
			HomeName:  "Berlin",
			AwayName:  "Düren",
			StartTime: 123,
			Country:   "Germany",
			League:    "Bundesliga Playoffs",
		}

		matchStore := db.NewMatchStore(MustOpenDB(t))
		importer := domain.NewMatchImporter(matchStore)

		_, err := importer.ImportMatches(ctx, untrackedMatch)
		expect.NoErr(t, err)

		matches, err := matchStore.All(ctx)
		expect.NoErr(t, err)

		expected := domain.Match{
			ID:        1,
			HomeName:  "Berlin",
			AwayName:  "Düren",
			StartTime: 123,
			Country:   "Germany",
			League:    "Bundesliga Playoffs",
		}
		expect.Equal(t, matches[0], expected)
	})

	//TODO test what happens if two matches with the same timestamp are in db
}
