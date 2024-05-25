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
		favs := []string{"Europe: Champions League - Play Offs", "USA: PVF Women"}

		err := importer.ImportMatches(ctx, favs)
		expect.NoErr(t, err)

		matches, err := matchStore.All(ctx)
		expect.NoErr(t, err)

		expected := domain.Matches{
			{
				HomeName:  "Trentino",
				AwayName:  "Jastrzebski",
				StartTime: 1714917600,
				Country:   "Europe",
				League:    "Champions League - Play Offs",
			},
			{
				HomeName:  "Resovia",
				AwayName:  "Zaksa",
				StartTime: 1714917600,
				Country:   "Europe",
				League:    "Champions League - Play Offs",
			},
			{
				HomeName:  "Grand Rapids Rise W",
				AwayName:  "San Diego Mojo W",
				StartTime: 1714939200,
				Country:   "USA",
				League:    "PVF Women",
			},
		}
		expect.MatchesEqual(t, matches, expected)
	})

	//TODO test what happens if two matches with the same timestamp are in db
	//TODO show errors when DB is not there
}
