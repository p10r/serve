package domain

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	NoFavouriteGamesTodayErr = errors.New("no favourite matches today")
	NoScheduledGamesTodayErr = errors.New("no scheduled matches today")
)

func (matches UntrackedMatches) FilterScheduled(favourites []string) (UntrackedMatches, error) {
	scheduled := filter("SCHEDULED", matches)
	if len(scheduled) == 0 {
		return nil, NoScheduledGamesTodayErr
	}

	filtered := filterFavourites(scheduled, favourites) //TODO
	if len(filtered) == 0 {
		return nil, NoFavouriteGamesTodayErr
	}

	return filtered, nil
}

func (matches UntrackedMatches) FilterFinished(favourites []string) (UntrackedMatches, error) {
	scheduled := filter("FINISHED", matches)
	if len(scheduled) == 0 {
		return nil, NoScheduledGamesTodayErr
	}

	filtered := filterFavourites(scheduled, favourites) //TODO
	if len(filtered) == 0 {
		return nil, NoFavouriteGamesTodayErr
	}

	return filtered, nil
}

func filter(stage string, flashscoreMatches UntrackedMatches) UntrackedMatches {
	filtered := UntrackedMatches{}

	for _, match := range flashscoreMatches {
		if lowerCase(match.Stage) == lowerCase(stage) {
			filtered = append(filtered, match)
		}
	}

	return filtered
}

// TODO move favourites to struct that has Country and League separate
func filterFavourites(matches UntrackedMatches, favourites []string) UntrackedMatches {
	var favs []string
	for _, favourite := range favourites {
		fmt.Printf("fav: %v\n", lowerCase(favourite))
		favs = append(favs, lowerCase(favourite))
	}

	filtered := UntrackedMatches{}
	for _, match := range matches {
		if slices.Contains(favs, lowerCase(match.FlashscoreName)) {
			filtered = append(filtered, match)
		}
	}

	return filtered
}

func lowerCase(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
