package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type Flashscore interface {
	GetUpcomingMatches() (UntrackedMatches, error)
}

type MatchImporter struct {
	store      MatchStore
	flashscore Flashscore
	favLeagues []string
}

func NewMatchImporter(store MatchStore, flashscore Flashscore, favLeagues []string) *MatchImporter {
	return &MatchImporter{store, flashscore, favLeagues}
}

// ImportScheduledMatches writes matches from flashscore into the db for the current day.
// Doesn't validate if the match is already present, as it's expected to be triggered only once per day for now.
func (importer MatchImporter) ImportScheduledMatches(ctx context.Context) error {
	untrackedMatches, err := importer.fetchAllMatches()
	if err != nil {
		return err
	}

	log.Printf("MatchImporter: %v matches coming up today", len(untrackedMatches))

	//TODO remove error, return empty slice
	upcoming, err := untrackedMatches.FilterScheduled(importer.favLeagues)
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

func (importer MatchImporter) fetchAllMatches() (UntrackedMatches, error) {
	untrackedMatches, err := importer.flashscore.GetUpcomingMatches()
	if err != nil {
		return nil, fmt.Errorf("could not fetch matches from flashscore: %v", err)
	}
	return untrackedMatches, err
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
