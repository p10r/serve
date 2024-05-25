package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type MatchImporter struct {
	store      MatchStore
	flashscore Flashscore
}

func NewMatchImporter(store MatchStore, flashscore Flashscore) *MatchImporter {
	return &MatchImporter{store, flashscore}
}

type Flashscore interface {
	GetUpcomingMatches() (UntrackedMatches, error)
}

type MatchStore interface {
	All(context.Context) (Matches, error)
	Add(context.Context, UntrackedMatch) (Match, error)
}

// ImportMatches writes matches from flashscore into the db for the current day.
// Doesn't validate if the match is already present, as it's expected to be triggered only once per day for now.
func (importer MatchImporter) ImportMatches(
	ctx context.Context,
	favLeagues []string,
) error {
	untrackedMatches, err := importer.flashscore.GetUpcomingMatches()
	if err != nil {
		return fmt.Errorf("could not fetch matches from flashscore: %v", err)
	}
	log.Printf("MatchImporter: %v matches upcoming today", len(untrackedMatches))

	//TODO remove error, return empty slice
	upcoming, err := untrackedMatches.FilterScheduled(favLeagues)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}

	_, err = importer.storeUntrackedMatch(ctx, upcoming)
	if err != nil {
		return err
	}

	return nil
}

func (importer MatchImporter) storeUntrackedMatch(ctx context.Context, matches UntrackedMatches) (Matches, error) {
	var trackedMatches Matches
	var dbErrs []error
	for _, untrackedMatch := range matches {
		trackedMatch, err := importer.store.Add(ctx, untrackedMatch)
		if err != nil {
			dbErr := fmt.Errorf("could not persist match %v, aborting: %v", untrackedMatch, err)
			dbErrs = append(dbErrs, dbErr)
		}

		trackedMatches = append(trackedMatches, trackedMatch)
	}

	return trackedMatches, errors.Join(dbErrs...)
}
