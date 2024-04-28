package domain

import (
	"errors"
	"github.com/p10r/serve/flashscore"
	"slices"
	"strings"
)

var (
	NoFavouriteGamesTodayErr = errors.New("no favourite matches today")
	NoScheduledGamesTodayErr = errors.New("no scheduled matches today")
)

func FilterScheduled(leagues flashscore.Leagues, favourites []string) (flashscore.Leagues, error) {
	scheduled := filter("Scheduled", leagues)
	if len(scheduled) == 0 {
		return nil, NoScheduledGamesTodayErr
	}

	filteredFavourites := filterFavourites(scheduled, favourites) //TODO
	if len(filteredFavourites) == 0 {
		return nil, NoFavouriteGamesTodayErr
	}

	return filteredFavourites, nil
}

func FilterFinished(leagues flashscore.Leagues, favourites []string) (flashscore.Leagues, error) {
	scheduled := filter("Finished", leagues)
	if len(scheduled) == 0 {
		return nil, NoScheduledGamesTodayErr
	}

	filteredFavourites := filterFavourites(scheduled, favourites) //TODO
	if len(filteredFavourites) == 0 {
		return nil, NoFavouriteGamesTodayErr
	}

	return filteredFavourites, nil
}

func filter(stage string, leagues flashscore.Leagues) flashscore.Leagues {
	filteredLeagues := flashscore.Leagues{}
	for _, league := range leagues {
		scheduledEvents := flashscore.Events{}
		for _, event := range league.Events {
			if sanitized(event.Stage) == sanitized(stage) {
				scheduledEvents = append(scheduledEvents, event)
			}
		}

		if len(scheduledEvents) > 0 {
			filteredLeagues = append(filteredLeagues, flashscore.League{
				Name:   league.Name,
				Events: scheduledEvents,
			})
		}
	}

	return filteredLeagues
}

func filterFavourites(leagues flashscore.Leagues, favourites []string) flashscore.Leagues {
	var sanitizedFavourites []string
	for _, favourite := range favourites {
		sanitizedFavourites = append(sanitizedFavourites, sanitized(favourite))
	}

	found := flashscore.Leagues{}
	for _, league := range leagues {
		if slices.Contains(sanitizedFavourites, sanitized(league.Name)) {
			found = append(found, league)
		}
	}
	return found
}

func sanitized(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
