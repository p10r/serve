package domain_test

import (
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/expect"
	"github.com/p10r/serve/testutil"
	"testing"
)

func TestDomain(t *testing.T) {
	t.Run("filters for scheduled matches", func(t *testing.T) {
		expected := domain.UntrackedMatches{
			{
				HomeName:       "Trentino",
				AwayName:       "Jastrzebski",
				StartTime:      1714917600,
				FlashscoreName: "Europe: Champions League - Play Offs",
				Country:        "Europe",
				League:         "Champions League - Play Offs",
				Stage:          "SCHEDULED",
			},
			{
				HomeName:       "Resovia",
				AwayName:       "Zaksa",
				StartTime:      1714917600,
				FlashscoreName: "Europe: Champions League - Play Offs",
				Country:        "Europe",
				League:         "Champions League - Play Offs",
				Stage:          "SCHEDULED",
			},
		}

		untrackedMatches := testutil.UntrackedMatches(t)
		matches, err := untrackedMatches.FilterScheduled([]string{"Europe: Champions League - Play Offs"})
		expect.NoErr(t, err)
		expect.DeepEqual(t, matches, expected)
	})

	t.Run("filters for finished matches", func(t *testing.T) {
		expected := domain.UntrackedMatches{
			{
				HomeName:       "Mok Mursa",
				AwayName:       "HAOK Mladost",
				StartTime:      1714932000,
				FlashscoreName: "Croatia: Superliga - Play Offs",
				Country:        "Croatia",
				League:         "Superliga - Play Offs",
				Stage:          "FINISHED",
			},
		}

		untrackedMatches := testutil.UntrackedMatches(t)
		matches, err := untrackedMatches.FilterFinished([]string{"Croatia: Superliga - Play Offs"})
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
			{
				HomeName:       "Trentino",
				AwayName:       "Jastrzebski",
				StartTime:      1714917600,
				FlashscoreName: "Europe: Champions League - Play Offs",
				Country:        "Europe",
				League:         "Champions League - Play Offs",
				Stage:          "SCHEDULED",
			},
			{
				HomeName:       "Resovia",
				AwayName:       "Zaksa",
				StartTime:      1714917600,
				FlashscoreName: "Europe: Champions League - Play Offs",
				Country:        "Europe",
				League:         "Champions League - Play Offs",
				Stage:          "SCHEDULED",
			},
			{
				HomeName:       "Grand Rapids Rise W",
				AwayName:       "San Diego Mojo W",
				StartTime:      1714939200,
				FlashscoreName: "USA: PVF Women",
				Country:        "USA",
				League:         "PVF Women",
				Stage:          "SCHEDULED",
			},
		}

		favourites := []string{"Europe: Champions League - Play Offs", "USA: PVF Women"}

		matches, err := testutil.UntrackedMatches(t).FilterScheduled(favourites)

		expect.NoErr(t, err)
		expect.DeepEqual(t, matches, expected)
	})
}
