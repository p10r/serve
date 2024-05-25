package domain

import "context"

type MatchStore interface {
	All(context.Context) (Matches, error)
	Add(context.Context, UntrackedMatch) (Match, error)
}
