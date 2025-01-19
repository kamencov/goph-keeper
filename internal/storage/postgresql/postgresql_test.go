package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log/slog"
	"os"
	"testing"
)

var (
	errCheckUser = errors.New("failed to check user")
	errCheckPass = errors.New("failed to check password")
	errSaveUser  = errors.New("failed to handlers user")
	errDelUser   = errors.New("failed to deleted credentials")
)

func TestPostgresql_Close(t *testing.T) {
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
	mock.ExpectClose()
	p := &Postgresql{
		log:     log,
		storage: db,
	}
	err = p.Close()
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

}

func TestPostgresql_CheckUser(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		testRow       *sqlmock.Rows
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_check",
			login:         "test",
			testRow:       sqlmock.NewRows([]string{"login"}).AddRow(1),
			expectedQuery: nil,
			expectedErr:   ErrUserAlreadyExists,
		},
		{
			name:          "failed_to_check_user",
			login:         "test",
			testRow:       sqlmock.NewRows([]string{"login"}),
			expectedQuery: errCheckUser,
			expectedErr:   errCheckUser,
		},
		{
			name:          "user_not_exists",
			login:         "test",
			testRow:       sqlmock.NewRows([]string{"login"}),
			expectedQuery: sql.ErrNoRows,
			expectedErr:   nil,
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
			mock.ExpectQuery("SELECT login").WillReturnRows(tt.testRow).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.CheckUser(context.Background(), tt.login)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_CheckPassword(t *testing.T) {
	tests := []struct {
		name             string
		login            string
		testRow          *sqlmock.Rows
		expectedQueryErr error
		expectedErr      bool
	}{
		{
			name:        "successful_check",
			login:       "test",
			testRow:     sqlmock.NewRows([]string{"password"}).AddRow(1),
			expectedErr: true,
		},
		{
			name:        "failed_to_check_password",
			login:       "test",
			testRow:     sqlmock.NewRows([]string{"password"}),
			expectedErr: false,
		},
		{
			name:             "failed_rows",
			login:            "test",
			testRow:          sqlmock.NewRows([]string{"password"}),
			expectedQueryErr: errCheckPass,
			expectedErr:      false,
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

			mock.ExpectQuery("SELECT password").WillReturnRows(tt.testRow).WillReturnError(tt.expectedQueryErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			_, ok := p.CheckPassword(tt.login)

			if ok != tt.expectedErr {
				t.Errorf("unexpected err: %v", ok)
			}
		})
	}
}

func TestPostgresql_SaveUser(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "user_already_exists",
			login:         "test",
			password:      "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "user_not_exists",
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
			mock.ExpectExec("INSERT INTO users").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.SaveUser(context.Background(), tt.login, tt.password)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_SaveTableUserAndUpdateToken(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		token         string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_save",
			login:         "test",
			token:         "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_save",
			login:         "test",
			token:         "test",
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
			mock.ExpectExec("UPDATE users SET").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.SaveTableUserAndUpdateToken(tt.login, tt.token)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_ServerSaveLoginAndPasswordInCredentials(t *testing.T) {
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
			mock.ExpectExec("INSERT INTO credentials").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.ServerSaveLoginAndPasswordInCredentials(context.Background(), tt.userID, tt.resource, tt.login, tt.password)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_SaveTextDataPstgres(t *testing.T) {
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
			mock.ExpectExec("INSERT INTO text_data").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.SaveTextDataPstgres(context.Background(), tt.userID, tt.data)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_SaveBinaryDataBinary(t *testing.T) {
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
			mock.ExpectExec("INSERT INTO binary_data").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.SaveBinaryDataBinary(context.Background(), tt.userID, tt.data)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_SaveCards(t *testing.T) {
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
			mock.ExpectExec("INSERT INTO cards").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.SaveCards(context.Background(), tt.userID, tt.data)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_GetUserIDByToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		testRow     *sqlmock.Rows
		expectedGet error
		expectedErr error
	}{
		{
			name:        "successful_get_user_id_by_token",
			token:       "test",
			testRow:     sqlmock.NewRows([]string{"token"}).AddRow(1),
			expectedErr: nil,
		},
		{
			name:        "failed_get_user_id_by_token",
			token:       "test",
			testRow:     sqlmock.NewRows([]string{"token"}),
			expectedGet: errCheckUser,
			expectedErr: errCheckUser,
		},
		{
			name:        "failed_rows",
			token:       "test",
			testRow:     sqlmock.NewRows([]string{"token"}).AddRow(1),
			expectedGet: sql.ErrNoRows,
			expectedErr: nil,
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
				WillReturnRows(tt.testRow).
				WillReturnError(tt.expectedGet)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			_, err = p.GetUserIDByToken(context.Background(), tt.token)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_GetUserIDByLogin(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		testRow     *sqlmock.Rows
		expectedDel error
		expectedErr error
	}{
		{
			name:        "successful_get_user_id_by_login",
			login:       "test",
			testRow:     sqlmock.NewRows([]string{"login"}).AddRow(1),
			expectedErr: nil,
		},
		{
			name:        "failed_get_user_id_by_login",
			login:       "test",
			testRow:     sqlmock.NewRows([]string{"login"}),
			expectedDel: errCheckUser,
			expectedErr: sql.ErrNoRows,
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

			mock.ExpectQuery("SELECT id").WillReturnRows(tt.testRow).WillReturnError(tt.expectedDel)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			_, err = p.GetUserIDByLogin(context.Background(), tt.login)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_DeletedCredentials(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		resource      string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_delete",
			userID:        1,
			resource:      "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_delete",
			userID:        1,
			resource:      "test",
			expectedQuery: errDelUser,
			expectedErr:   errDelUser,
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
			mock.ExpectExec("UPDATE credentials").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.DeletedCredentials(context.Background(), tt.userID, tt.resource)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_DeletedText(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		data          string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_delete",
			userID:        1,
			data:          "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_delete",
			userID:        1,
			data:          "test",
			expectedQuery: errDelUser,
			expectedErr:   errDelUser,
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
			mock.ExpectExec("UPDATE text_data").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.DeletedText(context.Background(), tt.userID, tt.data)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_DeletedBinary(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		data          string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_delete",
			userID:        1,
			data:          "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_delete",
			userID:        1,
			data:          "test",
			expectedQuery: errDelUser,
			expectedErr:   errDelUser,
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
			mock.ExpectExec("UPDATE binary_data").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.DeletedBinary(context.Background(), tt.userID, tt.data)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestPostgresql_DeletedCards(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		data          string
		expectedQuery error
		expectedErr   error
	}{
		{
			name:          "successful_delete",
			userID:        1,
			data:          "test",
			expectedQuery: nil,
			expectedErr:   nil,
		},
		{
			name:          "failed_delete",
			userID:        1,
			data:          "test",
			expectedQuery: errDelUser,
			expectedErr:   errDelUser,
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
			mock.ExpectExec("UPDATE cards").
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(tt.expectedErr)
			p := &Postgresql{
				log:     log,
				storage: db,
			}

			err = p.DeletedCards(context.Background(), tt.userID, tt.data)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}
