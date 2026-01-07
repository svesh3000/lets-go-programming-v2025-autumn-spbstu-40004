package db_test

import (
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
	t.Run("success with data", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		dbService := db.New(mockDB)

		testTable := [][]string{
			{"one", "two"},
			{"one", "two", "three"},
			{"one", ""},
			{"", ""},
		}

		for _, row := range testTable {
			mock.ExpectQuery("SELECT name FROM users").
				WillReturnRows(getMockDBRows(t, row))

			names, err := dbService.GetNames()

			require.NoError(t, err)
			require.Equal(t, row, names)
			require.NoError(t, mock.ExpectationsWereMet())
		}
	})
}
