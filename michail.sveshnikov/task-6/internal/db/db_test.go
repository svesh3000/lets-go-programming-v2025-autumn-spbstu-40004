package db_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/svesh3000/task-6/internal/db"
)

func getMockDBRows(t *testing.T, names []string) *sqlmock.Rows {
	t.Helper()

	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows = rows.AddRow(name)
	}

	return rows
}

func TestGetNames(t *testing.T) {
	successCases := []struct {
		name         string
		expectedRows []string
	}{
		{"two names", []string{"one", "two"}},
		{"three names", []string{"one", "two", "three"}},
		{"name with empty", []string{"one", ""}},
		{"empty names", []string{"", ""}},
	}

	for _, tc := range successCases {
		t.Run("success/"+tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			dbService := db.New(mockDB)

			mock.ExpectQuery("SELECT name FROM users").
				WillReturnRows(getMockDBRows(t, tc.expectedRows))

			names, err := dbService.GetNames()

			require.NoError(t, err)
			require.Equal(t, tc.expectedRows, names)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("database query error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		dbService := db.New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(fmt.Errorf("connection failed"))

		names, err := dbService.GetNames()

		require.Error(t, err)
		require.ErrorContains(t, err, "db query")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
