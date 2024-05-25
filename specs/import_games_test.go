package specifications_test

import (
	"context"
	"github.com/p10r/serve/db"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
	"github.com/p10r/serve/flashscore"
	"github.com/p10r/serve/testutil"
	"testing"
)

func TestImportMatches(t *testing.T) {
	ctx := context.TODO()

	t.Run("imports today's matches to db", func(t *testing.T) {
		apiKey := "apiKey"

		flashscoreServer := testutil.NewFlashscoreServer(t, apiKey)
		defer flashscoreServer.Close()

		client := flashscore.NewClient(flashscoreServer.URL, apiKey)
		matchStore := db.NewMatchStore(testutil.MustOpenDB(t))
		importer := domain.NewMatchImporter(matchStore, client)

		err := importer.ImportMatches(ctx)
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
	//TODO show errors when DB is not there
}
