package db_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/svesh3000/task-6/internal/db"
)

type testCase struct {
	name          string
	setupMock     func(sqlmock.Sqlmock)
	wantErr       bool
	errSubstr     string
	expectedNames []string
}

func createMockRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows = rows.AddRow(name)
	}

	return rows
}

func runDBTests(t *testing.T, testCases []testCase, method func(db.DBService) ([]string, error)) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			dbService := db.New(mockDB)
			tc.setupMock(mock)

			names, err := method(dbService)

			if tc.wantErr {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errSubstr)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetNames(t *testing.T) {
	testCases := []testCase{
		{
			name: "success: two names",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(createMockRows([]string{"one", "two"}))
			},
			wantErr:       false,
			expectedNames: []string{"one", "two"},
		},
		{
			name: "success: two empty names",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(createMockRows([]string{"", ""}))
			},
			wantErr:       false,
			expectedNames: []string{"", ""},
		},
		{
			name: "error: database query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(fmt.Errorf("connection failed"))
			},
			wantErr:   true,
			errSubstr: "db query",
		},
		{
			name: "error: row scanning error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(
						sqlmock.NewRows([]string{"name"}).
							AddRow(nil),
					)
			},
			wantErr:   true,
			errSubstr: "rows scanning",
		},
		{
			name: "error: rows error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(
						sqlmock.NewRows([]string{"name"}).
							AddRow("one").
							RowError(0, fmt.Errorf("row iteration failed")),
					)
			},
			wantErr:   true,
			errSubstr: "rows error",
		},
	}

	runDBTests(t, testCases, func(s db.DBService) ([]string, error) {
		return s.GetNames()
	})
}

func TestGetUniqueNames(t *testing.T) {
	testCases := []testCase{
		{
			name: "success: two names",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(createMockRows([]string{"one", "two"}))
			},
			wantErr:       false,
			expectedNames: []string{"one", "two"},
		},
		{
			name: "success: two empty names",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(createMockRows([]string{"", ""}))
			},
			wantErr:       false,
			expectedNames: []string{"", ""},
		},
		{
			name: "error: database query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(fmt.Errorf("connection failed"))
			},
			wantErr:   true,
			errSubstr: "db query",
		},
		{
			name: "error: row scanning error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(
						sqlmock.NewRows([]string{"name"}).
							AddRow(nil),
					)
			},
			wantErr:   true,
			errSubstr: "rows scanning",
		},
		{
			name: "error: rows error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(
						sqlmock.NewRows([]string{"name"}).
							AddRow("one").
							RowError(0, fmt.Errorf("row iteration failed")),
					)
			},
			wantErr:   true,
			errSubstr: "rows error",
		},
	}

	runDBTests(t, testCases, func(s db.DBService) ([]string, error) {
		return s.GetUniqueNames()
	})
}
