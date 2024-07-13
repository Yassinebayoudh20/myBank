package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)

	n := 5
	amount := int64(10)

	//!To get data from GoRoutines you need channels to await for the goroutine finish and assign those variables
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() { //! Its like creating a promise in ts
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	//? check Results
	for i := 0; i < n; i++ {
		err := <-errs //! Its like awaiting for the result from promise "errs" into err
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountFrom.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountTo.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountFrom.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountTo.ID, toAccount.ID)

		// check balances

		moneyGoingOutFromAccount := accountFrom.Balance - fromAccount.Balance
		moneyGoingInToAccount := toAccount.Balance - accountTo.Balance
		require.Equal(t, moneyGoingOutFromAccount, moneyGoingInToAccount)
		require.True(t, moneyGoingOutFromAccount > 0)
		require.True(t, moneyGoingOutFromAccount%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(moneyGoingOutFromAccount / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedFromAccount, err := store.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)

	require.Equal(t, accountFrom.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, accountTo.Balance+int64(n)*amount, updatedToAccount.Balance)

}
