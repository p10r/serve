package discord_test

import (
	"encoding/json"
	"github.com/p10r/serve/discord"
	"github.com/p10r/serve/flashscore"
	"github.com/p10r/serve/helpers"
	"testing"
)

func TestDiscord(t *testing.T) {
	leagues := flashscore.Leagues{
		{
			"Germany: Bundesliga",
			flashscore.Events{
				{
					HomeName:         "BR Volleys",
					AwayName:         "VfB Friedrichshafen",
					StartTime:        1697898600,
					HomeScoreCurrent: "3",
					AwayScoreCurrent: "1",
					Stage:            "FINISHED",
				},
			},
		},
		{
			"Italy: SuperLega",
			flashscore.Events{
				{
					HomeName:         "Lube",
					AwayName:         "Piacenza",
					StartTime:        1697898600,
					HomeScoreCurrent: "3",
					AwayScoreCurrent: "1",
					Stage:            "FINISHED",
				},
				{
					HomeName:         "Perugia",
					AwayName:         "Modena",
					StartTime:        1697965200,
					HomeScoreCurrent: "",
					AwayScoreCurrent: "",
					Stage:            "SCHEDULED",
				},
			},
		},
	}

	t.Run("creates discord message", func(t *testing.T) {
		msg := discord.NewMessage(leagues)
		marshal, err := json.MarshalIndent(msg, "", " ")
		helpers.NoErr(t, err)

		expected := helpers.ReadFile(t, "../helpers/expected-discord-message.json")

		helpers.JsonEqual(t, marshal, expected)
	})

	t.Run("sends discord Message", func(t *testing.T) {
		server := discord.NewFakeServer(t)
		defer server.Close()

		msg := discord.NewMessage(leagues)

		c := discord.NewClient(server.URL)
		err := c.SendMessage(msg)
		helpers.NoErr(t, err)
	})
}
