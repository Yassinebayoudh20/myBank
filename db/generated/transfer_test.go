package db

import (
	"context"
	"testing"

	"github.com/Yassinebayoudh20/my_bank/factory"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	accountFrom, accountTo := CreateTransferAccounts(t)
	createRandomTransfer(t, accountFrom, accountTo)
}

func CreateTransferAccounts(t *testing.T) (Account, Account) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	return accountFrom, accountTo
}

func TestGetListTransfers(t *testing.T) {
	accountFrom, accountTo := CreateTransferAccounts(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, accountFrom, accountTo)
	}

	params := ListTransfersParams{
		Limit:         5,
		Offset:        5,
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, entry := range transfers {
		require.NotEmpty(t, entry)
	}
}

func TestGetTransfer(t *testing.T) {
	accountFrom, accountTo := CreateTransferAccounts(t)
	transferCreate := createRandomTransfer(t, accountFrom, accountTo)
	transferResult, err := testQueries.GetTransfer(context.Background(), transferCreate.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferResult)

	require.Equal(t, transferCreate.ID, transferResult.ID)
	require.Equal(t, transferCreate.FromAccountID, transferResult.FromAccountID)
	require.Equal(t, transferCreate.ToAccountID, transferResult.ToAccountID)
	require.Equal(t, transferCreate.Amount, transferResult.Amount)

}

func createRandomTransfer(t *testing.T, accountFrom Account, accountTo Account) Transfer {

	params := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        factory.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, params.FromAccountID, transfer.FromAccountID)
	require.Equal(t, params.ToAccountID, transfer.ToAccountID)
	require.Equal(t, params.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}
