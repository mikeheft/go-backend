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

// TransferTxParams contains the input paramters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:from_account_id`
	ToAccountID   int64 `json:to_account_id`
	Amount        int64 `json:amount`
}

type TransferTxResult struct {
	Transfer    Transfer `json:transfer`
	FromAccount Account  `json:from_account`
	ToAccount   Account  `json:to_account`
	FromEntry   Entry    `json:from_entry`
	ToEntry     Entry    `json:to_entry`
}

// TransferTx performs a moeny transer from one account to the other.
// It creates a transer record, add account entries and update accounts' balance withing a db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		return nil
	})
}
