package specifications_test

import (
	"context"
	"encoding/json"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/p10r/serve/db"
	"github.com/p10r/serve/discord"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
	"github.com/p10r/serve/flashscore"
	"github.com/p10r/serve/testutil"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	r := approvals.UseReporter(reporters.NewIntelliJReporter())
	defer r.Close()

	approvals.UseFolder("testdata")
	os.Exit(m.Run())
}

type fixture struct {
	flashscoreServer *httptest.Server
	discordServer    *testutil.DiscordServer
	importer         *domain.MatchImporter
	store            *db.MatchStore
}

func newFixture(t *testing.T, favLeagues []string, runAgainstDiscord bool) fixture {
	apiKey := "random_api_key"
	flashscoreServer := testutil.NewFlashscoreServer(t, apiKey)
	fsClient := flashscore.NewClient(flashscoreServer.URL, apiKey)

	discordServer := testutil.NewDiscordServer(t)

	var discordClient *discord.Client
	if runAgainstDiscord {
		discordClient = discord.NewClient(os.Getenv("DISCORD_URI"))
	} else {
		discordClient = discord.NewClient(discordServer.URL)
	}

	may28th := func() time.Time {
		return time.Date(2024, 5, 28, 0, 0, 0, 0, time.UTC)
	}

	matchStore := db.NewMatchStore(testutil.MustOpenDB(t))
	importer := domain.NewMatchImporter(matchStore, fsClient, discordClient, favLeagues, may28th)
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

	f := newFixture(t, favs, false)
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
		requests := *f.discordServer.Requests
		expect.Len(t, requests, 1)

		var msg discord.Message
		err := json.Unmarshal(requests[0], &msg)
		expect.NoErr(t, err)
		msg = orderLeagues(msg)

		marshal, err := json.MarshalIndent(msg, "", " ")
		println(string(marshal))
		expect.NoErr(t, err)

		approvals.VerifyJSONBytes(t, marshal)
	})

	// Make sure to:
	// 1. remove t.Skip()
	// 2. direnv allow . && go test specs/scheduled_matches_test.go
	t.Run("run against real discord", func(t *testing.T) {
		//t.Skip()

		newFixture(t, favs, true).importer.ImportScheduledMatches(ctx)
	})

}

// we order the leagues to make sure the output json has always the same structure
func orderLeagues(msg discord.Message) discord.Message {
	sort.Slice(msg.Embeds[0].Fields, func(i, j int) bool {
		leagueName1 := msg.Embeds[0].Fields[i].Name
		leagueName2 := msg.Embeds[0].Fields[j].Name

		return len(leagueName1) < len(leagueName2)
	})

	return msg
}

//TODO test what happens if two matches with the same timestamp are in db
//TODO show errors when DB is not there
