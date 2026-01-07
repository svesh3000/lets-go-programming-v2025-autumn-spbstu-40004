package db_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/svesh3000/task-6/internal/db"
)

func TestGetNamesSuccess(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when creating mock db: %v", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	expectedNames := []string{"One", "Two"}

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(
			sqlmock.NewRows([]string{"numbers"}).
				AddRow("One").
				AddRow("Two"),
		)

	names, err := dbService.GetNames()

	require.NoError(t, err)
	require.Equal(t, expectedNames, names)

	require.NoError(t, mock.ExpectationsWereMet())
}
