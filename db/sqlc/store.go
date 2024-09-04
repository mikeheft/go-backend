package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execture db queries and transactions
type Store struct {
	// Composition > Inheritance
	// The preferred manner in which to extend a struct is to
	// by embedding
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// extecTx exectures a function within a transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %w, rb err: %w", err, rbErr)
		}
	}

	return tx.Commit()
}
