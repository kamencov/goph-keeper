package sqlite

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log/slog"
	"os"
	"testing"
)

var (
	errGetAllSync = errors.New("failed to get all sync")
	errGet        = errors.New("failed to get data")
	errClearSync  = errors.New("failed to clear sync")
)

func TestStorage_GetAllSync(t *testing.T) {
	tests := []struct {
		name             string
		testRows         *sqlmock.Rows
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name: "successful_get",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "table_name", "task_id", "action", "updated_at"}).
				AddRow(1, 1, "test", 1, "test", "2006-01-02 15:04:05"),
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_get_all_sync",
			testRows:         sqlmock.NewRows([]string{"id", "user_id", "table_name", "task_id", "action", "updated_at"}),
			expectedQueryErr: errGetAllSync,
			expectedErr:      errGetAllSync,
		},
		{
			name: "failed_rows",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "table_name", "task_id", "action", "updated_at"}).
				AddRow(1, 1, "test", 1, "test", "2006"),
			expectedQueryErr: nil,
			expectedErr:      errParseTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectQuery("SELECT id, user_id, table_name, task_id, action, updated_at").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetAllSync()
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetDataCredentials(t *testing.T) {
	tests := []struct {
		name             string
		userID           int
		taskID           int
		testRows         *sqlmock.Rows
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name: "successful_get",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "resource", "login", "password", "updated_at", "deleted_status"}).
				AddRow(1, 1, "test", "test", "test", "2006-01-02 15:04:05", 1),
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_get_data_credentials",
			testRows:         sqlmock.NewRows([]string{"id", "user_id", "resource", "login", "password", "updated_at", "deleted_status"}),
			expectedQueryErr: errGet,
			expectedErr:      errGet,
		},
		{
			name: "failed_parse_time",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "resource", "login", "password", "updated_at", "deleted_status"}).
				AddRow(1, 1, "test", "test", "test", "2006", 1),
			expectedQueryErr: nil,
			expectedErr:      errParseTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectQuery("SELECT *").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetDataCredentials(tt.userID, tt.taskID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetDataTextData(t *testing.T) {
	tests := []struct {
		name             string
		userID           int
		taskID           int
		testRows         *sqlmock.Rows
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name: "successful_get",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "text", "updated_at", "deleted_status"}).
				AddRow(1, 1, "test", "2006-01-02 15:04:05", 1),
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_get_data_text_data",
			testRows:         sqlmock.NewRows([]string{"id", "user_id", "text", "updated_at", "deleted_status"}),
			expectedQueryErr: errGet,
			expectedErr:      errGet,
		},
		{
			name: "failed_parse_time",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "text", "updated_at", "deleted_status"}).
				AddRow(1, 1, "test", "2006", 1),
			expectedQueryErr: nil,
			expectedErr:      errParseTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectQuery("SELECT *").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetDataTextData(tt.userID, tt.taskID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetDataBinaryData(t *testing.T) {
	tests := []struct {
		name             string
		userID           int
		taskID           int
		testRows         *sqlmock.Rows
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name: "successful_get",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "binary_data", "updated_at", "deleted_status"}).
				AddRow(1, 1, "test", "2006-01-02 15:04:05", 1),
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_get_data_text_data",
			testRows:         sqlmock.NewRows([]string{"id", "user_id", "binary_data", "updated_at", "deleted_status"}),
			expectedQueryErr: errGet,
			expectedErr:      errGet,
		},
		{
			name: "failed_parse_time",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "binary_data", "updated_at", "deleted_status"}).
				AddRow(1, 1, "test", "2006", 1),
			expectedQueryErr: nil,
			expectedErr:      errParseTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectQuery("SELECT *").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetDataBinaryData(tt.userID, tt.taskID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetDataCards(t *testing.T) {
	tests := []struct {
		name             string
		userID           int
		taskID           int
		testRows         *sqlmock.Rows
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name: "successful_get",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "cards", "updated_at", "deleted_status"}).
				AddRow(1, 1, "test", "2006-01-02 15:04:05", 1),
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_get_data_text_data",
			testRows:         sqlmock.NewRows([]string{"id", "cards", "binary_data", "updated_at", "deleted_status"}),
			expectedQueryErr: errGet,
			expectedErr:      errGet,
		},
		{
			name: "failed_parse_time",
			testRows: sqlmock.NewRows([]string{"id", "user_id", "cards", "updated_at", "deleted_status"}).
				AddRow(1, 1, "test", "2006", 1),
			expectedQueryErr: nil,
			expectedErr:      errParseTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectQuery("SELECT *").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetDataCards(tt.userID, tt.taskID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_ClearSyncCredentials(t *testing.T) {
	tests := []struct {
		name             string
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name:             "successful_clear",
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_clear",
			expectedQueryErr: errClearSync,
			expectedErr:      errClearSync,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.ClearSyncCredentials()
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_ClearSyncTextData(t *testing.T) {
	tests := []struct {
		name             string
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name:             "successful_clear",
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_clear",
			expectedQueryErr: errClearSync,
			expectedErr:      errClearSync,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.ClearSyncTextData()
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_ClearSyncBinaryData(t *testing.T) {
	tests := []struct {
		name             string
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name:             "successful_clear",
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_clear",
			expectedQueryErr: errClearSync,
			expectedErr:      errClearSync,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.ClearSyncBinaryData()
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_ClearSyncCards(t *testing.T) {
	tests := []struct {
		name             string
		expectedQueryErr error
		expectedErr      error
	}{
		{
			name:             "successful_clear",
			expectedQueryErr: nil,
			expectedErr:      nil,
		},
		{
			name:             "failed_clear",
			expectedQueryErr: errClearSync,
			expectedErr:      errClearSync,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectExec("DELETE FROM").WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tt.expectedQueryErr)

			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.ClearSyncCards()
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}
