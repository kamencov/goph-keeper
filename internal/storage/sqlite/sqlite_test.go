package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"log/slog"
	"os"
	"testing"
)

var (
	errCheckUser = errors.New("failed to check")
	errUpdate    = errors.New("failed to update")
	errSaveUser  = errors.New("failed to handlers user")
	errGetAll    = errors.New("failed to get all data from database")
	errDelete    = errors.New("failed to delete data from database")
)

func TestNewSqlStorage(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	db, err := NewSqlStorage(log)
	defer func(db *Storage) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	if err != nil {
		t.Errorf("NewSqlStorage error = %v", err)
	}

	// удаляем папку
	baseDir, err := os.Getwd()
	if err != nil {
		log.Error("Ошибка определения пути базы данных", "error", err)
		return
	}

	// Укажите имя каталога и файла базы данных
	dbPath := fmt.Sprintf("%s/storage", baseDir)

	err = os.RemoveAll(dbPath)
}

func TestStorage_GetUserIDWithLogin(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_user_id",
			login:         "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_to_check_user",
			login:         "test",
			testRows:      sqlmock.NewRows([]string{"id"}),
			expectedQuery: errCheckUser,
			expectedErr:   errCheckUser,
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
			mock.ExpectQuery("SELECT id").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetUserIDWithLogin(context.Background(), tt.login)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetUserIDWithToken(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_user_id",
			login:         "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_to_check_user",
			login:         "test",
			testRows:      sqlmock.NewRows([]string{"id"}),
			expectedQuery: errCheckUser,
			expectedErr:   errCheckUser,
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
			mock.ExpectQuery("SELECT id").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetUserIDWithToken(context.Background(), tt.login)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetUserPassword(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_user_password",
			login:         "test",
			testRows:      sqlmock.NewRows([]string{"password"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_to_check_password",
			login:         "test",
			testRows:      sqlmock.NewRows([]string{"password"}),
			expectedQuery: errCheckUser,
			expectedErr:   errCheckUser,
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
			mock.ExpectQuery("SELECT password").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetUserPassword(context.Background(), tt.login)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetUserToken(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_user_token",
			login:         "test",
			testRows:      sqlmock.NewRows([]string{"token"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_to_check_token",
			login:         "test",
			testRows:      sqlmock.NewRows([]string{"token"}),
			expectedQuery: errCheckUser,
			expectedErr:   errCheckUser,
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
			mock.ExpectQuery("SELECT token").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetUserToken(context.Background(), tt.login)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetTokenWithUserID(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_token",
			userID:        1,
			testRows:      sqlmock.NewRows([]string{"token"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_to_check_token",
			userID:        1,
			testRows:      sqlmock.NewRows([]string{"token"}),
			expectedQuery: errCheckUser,
			expectedErr:   errCheckUser,
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
			mock.ExpectQuery("SELECT token").
				WillReturnRows(tt.testRows).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetTokenWithUserID(context.Background(), tt.userID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_SaveLoginAndToken(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		token         string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_save",
			login:         "test",
			password:      "test",
			token:         "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_save",
			login:         "test",
			password:      "test",
			token:         "test",
			expectedQuery: errUpdate,
			expectedErr:   errUpdate,
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
			mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.SaveLoginAndToken(context.Background(), tt.login, tt.password, tt.token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_UpdateLoginAndToken(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		token         string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_update",
			userID:        1,
			token:         "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_update",
			userID:        1,
			token:         "test",
			expectedQuery: errUpdate,
			expectedErr:   errUpdate,
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
			mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.UpdateLoginAndToken(context.Background(), tt.userID, tt.token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_SaveLoginAndPasswordInCredentials(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		resource      string
		login         string
		password      string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_save",
			userID:        1,
			resource:      "test",
			login:         "test",
			password:      "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_save",
			userID:        1,
			resource:      "test",
			login:         "test",
			password:      "test",
			expectedQuery: errSaveUser,
			expectedErr:   errSaveUser,
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
			mock.ExpectExec("INSERT INTO credentials").WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.SaveLoginAndPasswordInCredentials(context.Background(), tt.userID, tt.resource, tt.login, tt.password)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_SaveTextDataInDatabase(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		data          string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_save",
			userID:        1,
			data:          "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_save",
			userID:        1,
			data:          "test",
			expectedQuery: errSaveUser,
			expectedErr:   errSaveUser,
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
			mock.ExpectExec("INSERT INTO text_data").WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.SaveTextDataInDatabase(context.Background(), tt.userID, tt.data)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_SaveBinaryDataInDatabase(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		data          string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_save",
			userID:        1,
			data:          "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_save",
			userID:        1,
			data:          "test",
			expectedQuery: errSaveUser,
			expectedErr:   errSaveUser,
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
			mock.ExpectExec("INSERT INTO binary_data").WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.SaveBinaryDataInDatabase(context.Background(), tt.userID, tt.data)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_SaveCardsInDatabase(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		data          string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_save",
			userID:        1,
			data:          "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_save",
			userID:        1,
			data:          "test",
			expectedQuery: errSaveUser,
			expectedErr:   errSaveUser,
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
			mock.ExpectExec("INSERT INTO cards").WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.SaveCardsInDatabase(context.Background(), tt.userID, tt.data)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_SaveSync(t *testing.T) {
	tests := []struct {
		name          string
		tableName     string
		userID        int
		taskID        int
		action        string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_save_sync",
			tableName:     "cards",
			userID:        1,
			taskID:        1,
			action:        "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_save_sync",
			tableName:     "cards",
			userID:        1,
			taskID:        1,
			action:        "test",
			expectedQuery: errSaveUser,
			expectedErr:   errSaveUser,
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
			mock.ExpectExec("INSERT INTO sync_client").WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.SaveSync(context.Background(), tt.tableName, tt.userID, tt.taskID, tt.action)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		tableName     string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_all",
			userID:        1,
			tableName:     "cards",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_get_all",
			userID:        1,
			tableName:     "cards",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: errGetAll,
			expectedErr:   errGetAll,
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
			mock.ExpectQuery("SELECT *").WillReturnRows(tt.testRows).WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetAll(context.Background(), tt.userID, tt.tableName)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetIDTaskCredentials(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		tableName     string
		task          string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_id_task_credentials",
			userID:        1,
			tableName:     "cards",
			task:          "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_get_id_task_credentials",
			userID:        1,
			tableName:     "cards",
			task:          "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: errGetAll,
			expectedErr:   errGetAll,
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
			mock.ExpectQuery("SELECT id").WillReturnRows(tt.testRows).WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetIDTaskCredentials(context.Background(), tt.tableName, tt.userID, tt.task)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetIDTaskText(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		tableName     string
		task          string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_id_task_text",
			userID:        1,
			tableName:     "cards",
			task:          "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_get_id_task_text",
			userID:        1,
			tableName:     "cards",
			task:          "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: errGetAll,
			expectedErr:   errGetAll,
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
			mock.ExpectQuery("SELECT id").WillReturnRows(tt.testRows).WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetIDTaskText(context.Background(), tt.tableName, tt.userID, tt.task)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetIDTaskBinary(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		tableName     string
		task          string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_id_task_binary",
			userID:        1,
			tableName:     "cards",
			task:          "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_get_id_task_binary",
			userID:        1,
			tableName:     "cards",
			task:          "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: errGetAll,
			expectedErr:   errGetAll,
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
			mock.ExpectQuery("SELECT id").WillReturnRows(tt.testRows).WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetIDTaskBinary(context.Background(), tt.tableName, tt.userID, tt.task)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_GetIDTaskCards(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		tableName     string
		task          string
		testRows      *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_get_id_task_cards",
			userID:        1,
			tableName:     "cards",
			task:          "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_get_id_task_cards",
			userID:        1,
			tableName:     "cards",
			task:          "test",
			testRows:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			expectedQuery: errGetAll,
			expectedErr:   errGetAll,
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
			mock.ExpectQuery("SELECT id").WillReturnRows(tt.testRows).WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			_, err = p.GetIDTaskCards(context.Background(), tt.tableName, tt.userID, tt.task)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestStorage_Deleted(t *testing.T) {
	tests := []struct {
		name          string
		tableName     string
		id            int
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_deleted",
			tableName:     "cards",
			id:            1,
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_deleted",
			tableName:     "cards",
			id:            1,
			expectedQuery: errDelete,
			expectedErr:   errDelete,
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
			mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.expectedQuery)
			p := &Storage{
				log:     log,
				storage: db,
			}
			err = p.Deleted(context.Background(), tt.tableName, tt.id)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}
