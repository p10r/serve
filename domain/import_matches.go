package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Flashscore interface {
	GetUpcomingMatches() (UntrackedMatches, error)
}

type Discord interface {
	SendMessage(context.Context, Matches, time.Time) error
}

type MatchImporter struct {
	store      MatchStore
	flashscore Flashscore
	discord    Discord
	favLeagues []string
	clock      func() time.Time
}

func NewMatchImporter(
	store MatchStore,
	flashscore Flashscore,
	discord Discord,
	favLeagues []string,
	clock func() time.Time,
) *MatchImporter {
	return &MatchImporter{store, flashscore, discord, favLeagues, clock}
}

// ImportScheduledMatches writes matches from flashscore into the db for the current day.
// Doesn't validate if the match is already present, as it's expected to be triggered only once per day for now.
func (importer MatchImporter) ImportScheduledMatches(ctx context.Context) (Matches, error) {
	untrackedMatches, err := importer.fetchAllMatches()
	if err != nil {
		return nil, err
	}

	log.Printf("MatchImporter: %v matches coming up today", len(untrackedMatches))

	//TODO remove error, return empty slice
	upcoming, err := untrackedMatches.FilterScheduled(importer.favLeagues)
	if err != nil {
		log.Printf("%v", err)
		return Matches{}, nil
	}

	trackedMatches, err := importer.storeUntrackedMatches(ctx, upcoming)
	if err != nil {
		return nil, err
	}

	err = importer.discord.SendMessage(ctx, trackedMatches, importer.clock())
	if err != nil {
		return nil, err
	}

	return trackedMatches, nil
}

func (importer MatchImporter) fetchAllMatches() (UntrackedMatches, error) {
	untrackedMatches, err := importer.flashscore.GetUpcomingMatches()
	if err != nil {
		return nil, fmt.Errorf("could not fetch matches from flashscore: %v", err)
	}
	return untrackedMatches, err
}

func (importer MatchImporter) storeUntrackedMatches(ctx context.Context, matches UntrackedMatches) (Matches, error) {
	var trackedMatches Matches
	var dbErrs []error
	for _, untrackedMatch := range matches {
		trackedMatch, err := importer.store.Add(ctx, untrackedMatch)
		if err != nil {
			dbErr := fmt.Errorf("could not persist match %v, aborting: %v", untrackedMatch, err)
			dbErrs = append(dbErrs, dbErr)
		}

		log.Println("Stored in DB: ", trackedMatch)

		trackedMatches = append(trackedMatches, trackedMatch)
	}

	return trackedMatches, errors.Join(dbErrs...)
}
