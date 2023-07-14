package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	var n = 5
	var amount = int64(10)
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)
	for i := 0; i < n; i++ {
		go func() {

			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        10,
			})
			errs <- err
			results <- res
		}()
	}
	for i := 0; i < n; i++ {
		chErr := <-errs
		res := <-results
		require.NoError(t, chErr)
		require.NotEmpty(t, res)
		// account test
		require.Equal(t, res.FromAccount.ID, account1.ID)
		require.Equal(t, res.ToAccount.ID, account2.ID)

		// transfer test
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err := store.GetTransfer(context.Background(), res.Transfer.ID)
		require.NoError(t, err)

		// entry test
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), res.FromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), res.ToEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		fmt.Println(">> tx: ", res.FromAccount.Balance, res.ToAccount.Balance)
		// check balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, 4 * amount
	}
	// updated balance
	fromUpdatedAccount, _ := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, fromUpdatedAccount)
	require.Equal(t, account1.Balance-int64(amount)*int64(n), fromUpdatedAccount.Balance)

	toUpdatedAccount, _ := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, toUpdatedAccount)
	require.Equal(t, account2.Balance+int64(amount)*int64(n), toUpdatedAccount.Balance)
	fmt.Println(">> after: ", fromUpdatedAccount.Balance, toUpdatedAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	var n = 10
	var amount = int64(10)
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
			results <- res
		}()
	}
	for i := 0; i < n; i++ {
		chErr := <-errs
		res := <-results
		require.NoError(t, chErr)
		require.NotEmpty(t, res)
	}
	// updated balance
	fromUpdatedAccount, _ := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, fromUpdatedAccount)
	require.Equal(t, account1.Balance, fromUpdatedAccount.Balance)

	toUpdatedAccount, _ := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, toUpdatedAccount)
	require.Equal(t, account2.Balance, toUpdatedAccount.Balance)
	fmt.Println(">> after: ", fromUpdatedAccount.Balance, toUpdatedAccount.Balance)
}
