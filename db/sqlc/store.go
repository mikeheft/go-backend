package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier

	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execture SQL queries and transactions
type SQLStore struct {
	// Composition > Inheritance
	// The preferred manner in which to extend a struct is to
	// by embedding
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// extecTx exectures a function within a transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		// Rollback the transaction if any error occurs
		if rbErr := tx.Rollback(); rbErr != nil {
			// Return rollback error if rollback itself fails
			return fmt.Errorf("tx error: %w, rb err: %w", err, rbErr)
		}
		// Return the original transaction error after successful rollback
		return err
	}

	return tx.Commit()
}
