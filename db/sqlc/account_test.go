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
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	_, err = testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
}

func TestGetAccountByID(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
}
func TestGetAccountForUpdate(t *testing.T) {
	TestGetAccountByID(t)
}

func TestListAccounts(t *testing.T) {
	accounts := make([]Account, 5)
	for i := 0; i < 5; i++ {
		accounts[i] = CreateRandomAccount(t)
	}
	lap := ListAccountParams{
		Limit:  5,
		Offset: 0,
	}
	res, err := testQueries.ListAccount(context.Background(), lap)
	require.NoError(t, err)
	require.Equal(t, 5, len(res))
}

func TestUpdateAccounts(t *testing.T) {
	account := CreateRandomAccount(t)
	account.Balance = utils.RandomNumber(false)
	uap := UpdateAccountParams{
		ID:      account.ID,
		Balance: utils.RandomNumber(false),
	}
	res, err := testQueries.UpdateAccount(context.Background(), uap)
	require.NoError(t, err)
	require.Equal(t, res.ID, account.ID)

}
