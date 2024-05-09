package discord_test

import (
	"encoding/json"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/p10r/serve/discord"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
	"os"
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
	leagues := domain.UntrackedMatches{
		{
			HomeName:       "BR Volleys",
			AwayName:       "VfB Friedrichshafen",
			StartTime:      0,
			FlashscoreName: "Germany: 1. Bundesliga",
			Country:        "Germany",
			League:         "1. Bundesliga",
			Stage:          "FINISHED",
		},
		{
			HomeName:       "Lube",
			AwayName:       "Piacenza",
			StartTime:      1697898600,
			FlashscoreName: "Italy: SuperLega",
			Country:        "Italy",
			League:         "SuperLega",
			Stage:          "FINISHED",
		},
		{
			HomeName:       "Perugia",
			AwayName:       "Modena",
			StartTime:      1697965200,
			FlashscoreName: "Italy: SuperLega",
			Country:        "Italy",
			League:         "SuperLega",
			Stage:          "FINISHED",
		},
	}
	may28th := time.Date(2024, 5, 28, 0, 0, 0, 0, time.UTC)

	t.Run("creates discord message", func(t *testing.T) {
		msg := discord.NewMessage(leagues, may28th)
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
		server := discord.NewFakeServer(t)
		defer server.Close()

		msg := discord.NewMessage(leagues, may28th)

		c := discord.NewClient(server.URL)
		err := c.SendMessage(msg)
		expect.NoErr(t, err)
	})
}
