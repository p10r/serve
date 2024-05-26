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

type fixture struct {
	flashscoreServer *httptest.Server
	discordServer    *testutil.DiscordServer
	importer         *domain.MatchImporter
	store            *db.MatchStore
}

func newFixture(t *testing.T, favLeagues []string) fixture {
	apiKey := "random_api_key"
	flashscoreServer := testutil.NewFlashscoreServer(t, apiKey)
	fsClient := flashscore.NewClient(flashscoreServer.URL, apiKey)

	discordServer := testutil.NewDiscordServer(t)
	//discordClient := discord.NewClient(discordServer.URL)

	matchStore := db.NewMatchStore(testutil.MustOpenDB(t))
	importer := domain.NewMatchImporter(matchStore, fsClient, favLeagues)

	return fixture{
		flashscoreServer,
		discordServer,
		importer,
		matchStore,
	}
}

func TestImportMatches(t *testing.T) {
	ctx := context.TODO()
	favs := []string{"Europe: Champions League - Play Offs", "USA: PVF Women"}

	f := newFixture(t, favs)
	defer f.flashscoreServer.Close()
	defer f.discordServer.Close()

	_, err := f.importer.ImportScheduledMatches(ctx)
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
		expect.MatchStoreContains(t, f.store, expected)
	})

	t.Run("sends discord message", func(t *testing.T) {
		expect.Len(t, f.discordServer.Requests, 1)
	})

	//TODO test what happens if two matches with the same timestamp are in db
	//TODO show errors when DB is not there
}
