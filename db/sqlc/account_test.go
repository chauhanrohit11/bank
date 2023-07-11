package db

import (
	"context"
	"strings"
	"testing"

	"github.com/chauhanrohit11/bank/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	accountObj := CreateAccountParams{
		Owner:    utils.RandomString(10),
		Balance:  utils.RandomNumber(false),
		Currency: strings.ToUpper(utils.RandomString(3)),
	}
	account, err := testQueries.CreateAccount(context.Background(), accountObj)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.Owner, accountObj.Owner)
	require.Equal(t, account.Currency, accountObj.Currency)
	require.Equal(t, account.Balance, accountObj.Balance)
	require.NotZero(t, account.ID)
	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestDeleteAccount(t *testing.T) {
	account := CreateRandomAccount(t)
	deletedAccount := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NotEmpty(t, deletedAccount)
}

func TestGetAccountByID(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
}

func TestListAccounts(t *testing.T) {

}

func TestUpdateAccounts(t *testing.T) {}
