package specifications_test

import (
	"context"
	"github.com/p10r/serve/db"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
	"github.com/p10r/serve/flashscore"
	"github.com/p10r/serve/testutil"
	"net/http/httptest"
	"testing"
)

func newFixture(t *testing.T, favLeagues []string) (*db.MatchStore, *httptest.Server, *domain.MatchImporter) {
	apiKey := "random_api_key"
	flashscoreServer := testutil.NewFlashscoreServer(t, apiKey)

	client := flashscore.NewClient(flashscoreServer.URL, apiKey)
	matchStore := db.NewMatchStore(testutil.MustOpenDB(t))
	importer := domain.NewMatchImporter(matchStore, client, favLeagues)

	return matchStore, flashscoreServer, importer
}

func TestImportMatches(t *testing.T) {
	ctx := context.TODO()
	favs := []string{"Europe: Champions League - Play Offs", "USA: PVF Women"}

	matchStore, flashscoreServer, importer := newFixture(t, favs)
	defer flashscoreServer.Close()

	err := importer.ImportScheduledMatches(ctx)
	expect.NoErr(t, err)

	t.Run("imports today's matches to db", func(t *testing.T) {
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
		expect.MatchStoreContains(t, matchStore, expected)
	})

	t.Run("sends discord message", func(t *testing.T) {

	})

	//TODO test what happens if two matches with the same timestamp are in db
	//TODO show errors when DB is not there
}
