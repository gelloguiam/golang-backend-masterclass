package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

/* Transfer TX do the following:
1. Create transfer data from account 1 to account 2
2. Create entry data for account 1
3. Create entry data for account 2
4. Update balance of account 1
5. Update balance of account 2
*/
func (store *Store) transferTx(ctx context.Context, transferTxParams TransferTxParams) (result TransferTxResult, err error) {
	err = store.execTx(ctx, func(q *Queries) error {
		transfer, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: transferTxParams.FromAccountID,
			ToAccountID:   transferTxParams.ToAccountID,
			Amount:        transferTxParams.Amount,
		})

		if err != nil {
			return err
		}

		fromEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: transferTxParams.FromAccountID,
			Amount:    -transferTxParams.Amount,
		})
		if err != nil {
			return err
		}

		toEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: transferTxParams.ToAccountID,
			Amount:    transferTxParams.Amount,
		})
		if err != nil {
			return err
		}

		fromAccount, err := q.GetAccountForUpdate(ctx, transferTxParams.FromAccountID)
		if err != nil {
			return err
		}

		if fromAccount.Balance < transferTxParams.Amount {
			return fmt.Errorf("Insufficient balance")
		}

		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      transferTxParams.FromAccountID,
			Balance: fromAccount.Balance - transferTxParams.Amount,
		})
		if err != nil {
			return err
		}

		fromAccount, err = q.GetAccount(ctx, transferTxParams.FromAccountID)
		if err != nil {
			return err
		}

		toAccount, err := q.GetAccountForUpdate(ctx, transferTxParams.ToAccountID)
		if err != nil {
			return err
		}

		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      transfer.ToAccountID,
			Balance: toAccount.Balance + transferTxParams.Amount,
		})
		if err != nil {
			return err
		}

		toAccount, err = q.GetAccount(ctx, transferTxParams.ToAccountID)
		if err != nil {
			return err
		}

		result = TransferTxResult{
			Transfer:    transfer,
			FromAccount: fromAccount,
			ToAccount:   toAccount,
			FromEntry:   fromEntry,
			ToEntry:     toEntry,
		}

		return nil
	})

	return
}
