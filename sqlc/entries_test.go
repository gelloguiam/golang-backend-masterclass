package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	retrievedEntries, err := testQueries.GetEntriesByAccountID(context.Background(), entry.AccountID)

	require.NoError(t, err)
	require.NotEmpty(t, retrievedEntries)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	retrievedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, retrievedEntry)
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID:     entry.ID,
		Amount: util.RandomMoney(),
	}

	err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)

	retrievedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedEntry)
}
