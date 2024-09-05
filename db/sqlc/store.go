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

type TransferTxResult struct {
	Transfer    Transfer `json:transfer`
	FromAccount Account  `json:from_account`
	ToAccount   Account  `json:to_account`
	FromEntry   Entry    `json:from_entry`
	ToEntry     Entry    `json:to_entry`
}

var txKey = struct{}{}

// TransferTx performs a moeny transer from one account to the other.
// It creates a transer record, add account entries and update accounts' balance withing a db transaction
func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get from_account")
		// get account -> update its balance
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "update from_account")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get to_account")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		fmt.Println(txName, "update to_account")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
