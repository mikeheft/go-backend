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

type TransferTxResult struct {
	Transfer    Transfer `json:transfer`
	FromAccount Account  `json:from_account`
	ToAccount   Account  `json:to_account`
	FromEntry   Entry    `json:from_entry`
	ToEntry     Entry    `json:to_entry`
}

// TransferTx performs a moeny transer from one account to the other.
// It creates a transer record, add account entries and update accounts' balance withing a db transaction
func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create entry for fromAccount (debit)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create entry for toAccount (credit)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update from_account (debit)
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update to_account (credit)
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
