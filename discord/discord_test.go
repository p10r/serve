package discord_test

import (
	"encoding/json"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/p10r/serve/discord"
	"github.com/p10r/serve/expect"
	"github.com/p10r/serve/testutil"
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

func TestDiscord(t *testing.T) {
	leagues := testutil.UntrackedMatches(t)
	may28th := time.Date(2024, 5, 28, 0, 0, 0, 0, time.UTC)

	t.Run("creates discord message", func(t *testing.T) {
		initial := discord.NewMessage(leagues, may28th)

		// we order the leagues, otherwise it might differ each test run
		msg := orderLeagues(initial)

		marshal, err := json.MarshalIndent(msg, "", " ")
		expect.NoErr(t, err)

		approvals.VerifyJSONBytes(t, marshal)
	})

	// Make sure to:
	// 1. run direnv allow .
	// 2. run from command line
	t.Run("manually check", func(t *testing.T) {
		t.Skip()

		uri := os.Getenv("DISCORD_URI")
		println(uri)
		client := discord.NewClient(uri)

		err := client.SendMessage(discord.NewMessage(leagues, time.Now()))
		expect.NoErr(t, err)
	})

	t.Run("sends discord Message", func(t *testing.T) {
		server := testutil.NewDiscordServer(t)
		defer server.Close()

		msg := discord.NewMessage(leagues, may28th)

		c := discord.NewClient(server.URL)
		err := c.SendMessage(msg)
		expect.NoErr(t, err)
	})
}

func orderLeagues(msg discord.Message) discord.Message {
	sort.Slice(msg.Embeds[0].Fields, func(i, j int) bool {
		leagueName1 := msg.Embeds[0].Fields[i].Name
		leagueName2 := msg.Embeds[0].Fields[j].Name

		return len(leagueName1) < len(leagueName2)
	})

	return msg
}
