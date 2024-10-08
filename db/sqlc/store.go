package db

import (
    "context"
    "database/sql"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB)*Store {
	return &Store {
		db: db,
		Queries: New(db),
	}
}

func (store *Store)execTx(ctx context,Context, fn func(*Queries) error) error{
	tx,err := store.db.BeginTx(ctx,nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
    FromAccountID int64  `json:"from_account_id"`  // ID of the account to transfer from
    ToAccountID   int64  `json:"to_account_id"`    // ID of the account to transfer to
    Amount        int64  `json:"amount"`           // Transfer amount
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
    Transfer   Transfer `json:"transfer"`        // The transfer details
    FromAccount Account  `json:"from_account"`   // The account from which funds are transferred
    ToAccount   Account  `json:"to_account"`     // The account to which funds are transferred
    FromEntry   Entry    `json:"from_entry"`     // Entry for the source account
    ToEntry     Entry    `json:"to_entry"`       // Entry for the destination account
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	txName := ctx.Value(txKey)

	err := store.exeTx(ctx, func(q *Queries) error{
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: 	arg.FromAccountID,
			ToAccountID:	arg.ToAccountID,
			Amount:			arg.Amount,
		})
		if err != nil {
			return err
		}
		
		fmt.Println(exName,"create transfer")

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // Negative amount for the account transferring out
		})
		if err != nil {
			return err
		}
		
		fmt.Println(exName,"create entry 1")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // Positive amount for the account receiving funds
		})
		if err != nil {
			return err
		}

		fmt.Println(exName,"create entry 2")

		result.FromAccount, err = q.GetUpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}
		
		fmt.Println(exName,"get account 1")

		account2, err := q.GetUpdateAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		
		fmt.Println(exName,"create entry 1")

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