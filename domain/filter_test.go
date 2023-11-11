package domain_test

import (
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/flashscore"
	"github.com/p10r/serve/helpers"
	"testing"
)

func TestDomain(t *testing.T) {
	leagues := flashscore.Leagues{
		{
			"Austria: AVL",
			flashscore.Events{
				{
					HomeName:         "Sokol Wien",
					AwayName:         "VBK Klagenfurt",
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
		{
			"Latvia: AVL",
			flashscore.Events{
				{
					HomeName:         "Riga",
					AwayName:         "Jelgava",
					StartTime:        1697898600,
					HomeScoreCurrent: "3",
					AwayScoreCurrent: "1",
					Stage:            "SCHEDULED",
				},
			},
		},
	}

	t.Run("filters for scheduled games", func(t *testing.T) {
		expected := flashscore.Leagues{
			{
				"Italy: SuperLega",
				flashscore.Events{
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
			{
				"Latvia: AVL",
				flashscore.Events{
					{
						HomeName:         "Riga",
						AwayName:         "Jelgava",
						StartTime:        1697898600,
						HomeScoreCurrent: "3",
						AwayScoreCurrent: "1",
						Stage:            "SCHEDULED",
					},
				},
			},
		}

		leagues, err := domain.Filter(leagues, []string{"Italy: SuperLega", "Latvia: AVL"})
		helpers.NoErr(t, err)
		helpers.DeepEqual(t, leagues, expected)

	})

	t.Run("filters for favourites", func(t *testing.T) {
		expected := flashscore.Leagues{
			{
				"Italy: SuperLega",
				flashscore.Events{
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

		favourites := []string{"Italy: SuperLega"}

		leagues, err := domain.Filter(leagues, favourites)
		helpers.NoErr(t, err)
		helpers.DeepEqual(t, leagues, expected)
	})
}
