package db

import (
	"context"
	"testing"

	"github.com/Yassinebayoudh20/my_bank/factory"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	params := ListEntriesParams{
		Limit:     5,
		Offset:    5,
		AccountID: account.ID,
	}

	entries, err := testQueries.ListEntries(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entryCreate := createRandomEntry(t, account)
	entryResult, err := testQueries.GetEntry(context.Background(), entryCreate.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryResult)

	require.Equal(t, entryCreate.ID, entryResult.ID)
	require.Equal(t, entryCreate.AccountID, entryResult.AccountID)
	require.Equal(t, entryCreate.Amount, entryResult.Amount)

}

func createRandomEntry(t *testing.T, account Account) Entry {

	params := CreateEntryParams{
		AccountID: account.ID,
		Amount:    factory.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}
