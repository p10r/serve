package domain_test

import (
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
	"testing"
)

func TestDomain(t *testing.T) {
	//TODO move to one JSON file
	upcomingMatches := domain.UntrackedMatches{
		{
			HomeName:       "Sokol Wien",
			AwayName:       "VBK Klagenfurt",
			StartTime:      1697898600,
			FlashscoreName: "Austria: AVL",
			Country:        "Austria",
			League:         "AVL",
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
			Stage:          "SCHEDULED",
		},
		{
			HomeName:       "Riga",
			AwayName:       "Jelgava",
			StartTime:      1697898600,
			FlashscoreName: "Latvia: AVL",
			Country:        "Latvia",
			League:         "AVL",
			Stage:          "SCHEDULED",
		},
	}

	t.Run("filters for scheduled matches", func(t *testing.T) {
		expected := domain.UntrackedMatches{
			{
				HomeName:       "Perugia",
				AwayName:       "Modena",
				StartTime:      1697965200,
				FlashscoreName: "Italy: SuperLega",
				Country:        "Italy",
				League:         "SuperLega",
				Stage:          "SCHEDULED",
			},
			{
				HomeName:       "Riga",
				AwayName:       "Jelgava",
				StartTime:      1697898600,
				FlashscoreName: "Latvia: AVL",
				Country:        "Latvia",
				League:         "AVL",
				Stage:          "SCHEDULED",
			},
		}

		matches, err := upcomingMatches.FilterScheduled([]string{"Italy: SuperLega", "Latvia: AVL"})
		expect.NoErr(t, err)
		expect.DeepEqual(t, matches, expected)
	})

	t.Run("filters for finished matches", func(t *testing.T) {
		expected := domain.UntrackedMatches{
			{
				HomeName:       "Lube",
				AwayName:       "Piacenza",
				StartTime:      1697898600,
				FlashscoreName: "Italy: SuperLega",
				Country:        "Italy",
				League:         "SuperLega",
				Stage:          "FINISHED",
			},
		}

		matches, err := upcomingMatches.FilterFinished([]string{"Italy: SuperLega", "Latvia: AVL"})
		expect.NoErr(t, err)
		expect.DeepEqual(t, matches, expected)
	})

	t.Run("handles 0 scheduled matches", func(t *testing.T) {
		_, err := domain.UntrackedMatches{}.FilterScheduled([]string{"Italy: SuperLega"})
		expect.Err(t, err)
		expect.DeepEqual(t, err, domain.NoScheduledGamesTodayErr)
	})

	t.Run("filters for favourites", func(t *testing.T) {
		expected := domain.UntrackedMatches{
			domain.UntrackedMatch{
				HomeName:       "Perugia",
				AwayName:       "Modena",
				StartTime:      1697965200,
				FlashscoreName: "Italy: SuperLega",
				Country:        "Italy",
				League:         "SuperLega",
				Stage:          "SCHEDULED",
			},
		}

		favourites := []string{"Italy: SuperLega"}

		matches, err := upcomingMatches.FilterScheduled(favourites)
		expect.NoErr(t, err)
		expect.DeepEqual(t, matches, expected)
	})
}
