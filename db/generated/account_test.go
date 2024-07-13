package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Yassinebayoudh20/my_bank/factory"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	//create account
	accountCreate := createRandomAccount(t)
	accountResult, err := testQueries.GetAccount(context.Background(), accountCreate.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountResult)

	require.Equal(t, accountCreate.ID, accountResult.ID)
	require.Equal(t, accountCreate.Owner, accountResult.Owner)
	require.Equal(t, accountCreate.Balance, accountResult.Balance)
	require.Equal(t, accountCreate.Currency, accountResult.Currency)

}

func TestUpdateAccount(t *testing.T) {
	accountCreate := createRandomAccount(t)

	params := UpdateAccountParams{
		ID:      accountCreate.ID,
		Balance: factory.RandomMoney(),
	}

	accountResult, err := testQueries.UpdateAccount(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, accountResult)

	require.Equal(t, accountCreate.ID, accountResult.ID)
	require.Equal(t, accountCreate.Owner, accountResult.Owner)
	require.Equal(t, params.Balance, accountResult.Balance)
	require.Equal(t, accountCreate.Currency, accountResult.Currency)
}

func TestDeleteAccount(t *testing.T) {
	accountCreate := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), accountCreate.ID)
	require.NoError(t, err)

	accountResult, err := testQueries.GetAccount(context.Background(), accountCreate.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountResult)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	params := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func createRandomAccount(t *testing.T) Account {
	params := CreateAccountParams{
		Owner:    factory.RandomOwner(),
		Balance:  factory.RandomMoney(),
		Currency: factory.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
