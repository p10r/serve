package db

import (
	"context"
	"fmt"
	"github.com/p10r/serve/domain"
)

type MatchStore struct {
	db *DB
}

func NewMatchStore(db *DB) *MatchStore {
	return &MatchStore{db}
}

func (s MatchStore) All(ctx context.Context) (domain.Matches, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Matches{}, err
	}
	//nolint:errcheck
	defer tx.Rollback()

	matches, err := findAll(ctx, tx)
	if err != nil {
		return domain.Matches{}, err
	}

	//TODO remove pointer right in findMatches?
	var m domain.Matches
	for _, match := range matches {
		m = append(m, *match)
	}

	return m, err
}

func (s MatchStore) Add(ctx context.Context, untrackedMatch domain.UntrackedMatch) (domain.Match, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Match{}, err
	}
	//nolint:errcheck
	defer tx.Rollback()

	match, err := insertNewMatch(ctx, tx, untrackedMatch)
	if err != nil {
		return domain.Match{}, fmt.Errorf("could not insert match: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Match{}, fmt.Errorf("could not commit tx %v", err)
	}

	return match, nil
}

func insertNewMatch(ctx context.Context, tx *Tx, um domain.UntrackedMatch) (domain.Match, error) {
	result, err := tx.ExecContext(ctx, `
		INSERT INTO matches (
			home_name,
			away_name,
			start_time,
			country,
			league
		) VALUES (?,?,?,?,?)`, um.HomeName, um.AwayName, um.StartTime, um.Country, um.League)
	if err != nil {
		return domain.Match{}, err
	}

	// Read back new um ID into caller argument.
	id, err := result.LastInsertId()
	if err != nil {
		return domain.Match{}, err
	}

	return domain.NewMatch(id, um), nil
}

func findAll(ctx context.Context, tx *Tx) ([]*domain.Match, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT 
			id,
			home_name,
			away_name,
			start_time,
			country,
			league
		FROM matches 
		ORDER BY start_time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]*domain.Match, 0)
	for rows.Next() {
		var e domain.Match
		err := rows.Scan(
			&e.ID,
			&e.HomeName,
			&e.AwayName,
			&e.StartTime,
			&e.Country,
			&e.League,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
