package domain

import (
	"context"
)

type MatchImporter struct {
	store MatchStore
}

func NewMatchImporter(store MatchStore) *MatchImporter {
	return &MatchImporter{store: store}
}

type Flashscore interface {
	GetUpcomingMatches()
}

type MatchStore interface {
	All(context.Context) (Matches, error)
	Add(context.Context, UntrackedMatch) (Match, error)
}

// ImportMatches writes matches from flashscore into the db for the current day.
// Doesn't validate if the match is already present, as it's expected to be triggered only once per day for now.
func (importer MatchImporter) ImportMatches(
	ctx context.Context,
	untrackedMatch UntrackedMatch,
) (Matches, error) {
	trackedMatch, err := importer.store.Add(ctx, untrackedMatch)
	if err != nil {
		return nil, err
	}

	return Matches{trackedMatch}, nil
}
