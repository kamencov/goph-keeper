package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("the user already exists")
)

type Storage struct {
	storage *sql.DB
	log     *slog.Logger
}

func NewSqlStorage(log *slog.Logger) (*Storage, error) {
	db := &Storage{
		log: log,
	}

	// Получаем путь для базы данных
	dbPath, err := db.getDatabaseFilePath()
	if err != nil {
		log.Error("Ошибка определения пути базы данных", err)
		return nil, err
	}

	// Создаём базу данных
	if err := db.init(dbPath); err != nil {
		log.Error("Ошибка создания базы данных", err)
		return nil, err
	}

	fmt.Println("База данных и таблица успешно созданы!")

	return db, err
}

func (s *Storage) init(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	s.storage = db
	s.log.Info("connected to database")

	err = s.createTableIfNotExists()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) getDatabaseFilePath() (string, error) {
	baseDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Укажите имя каталога и файла базы данных
	dbPath := fmt.Sprintf("%s/storage", baseDir)
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/client.db", dbPath), nil
}

// createTableIfNotExists - создает таблицы в базе данных, если они не существуют.
func (s *Storage) createTableIfNotExists() (err error) {
	// Открытие транзакции
	tx, err := s.storage.Begin()
	if err != nil {
		s.log.Error("failed to begin transaction:", "error", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback() // Rollback, если Commit не был вызван
		}
	}()

	// Создаем таблицу users
	query := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        login TEXT NOT NULL UNIQUE,
        token TEXT,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    )`
	_, err = tx.Exec(query)
	if err != nil {
		s.log.Error("failed to create table - users:", "error", err)
		return err
	}

	// Создаем таблицу credentials
	query = `CREATE TABLE IF NOT EXISTS credentials (
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        user_id INTEGER NOT NULL, 
        resource TEXT NOT NULL,
        login TEXT NOT NULL,
        password TEXT NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`
	_, err = tx.Exec(query)
	if err != nil {
		s.log.Error("failed to create table - credentials:", "error", err)
		return err
	}

	// Создаем таблицу text_data
	query = `CREATE TABLE IF NOT EXISTS text_data (
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        user_id INTEGER NOT NULL, 
        text TEXT NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`
	_, err = tx.Exec(query)
	if err != nil {
		s.log.Error("failed to create table - text_data:", "error", err)
		return err
	}

	// Создаем таблицу binary_data
	query = `CREATE TABLE IF NOT EXISTS binary_data (
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        user_id INTEGER NOT NULL, 
        binary_data BLOB NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`
	_, err = tx.Exec(query)
	if err != nil {
		s.log.Error("failed to create table - binary_data:", "error", err)
		return err
	}

	// Создаем таблицу cards
	query = `CREATE TABLE IF NOT EXISTS cards (
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        user_id INTEGER NOT NULL, 
        cards TEXT NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`
	_, err = tx.Exec(query)
	if err != nil {
		s.log.Error("failed to create table - cards:", "error", err)
		return err
	}

	// Подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		s.log.Error("failed to commit transaction:", "error", err)
		return err
	}

	return nil
}

// Close закрывает соединение с базой данных.
func (s *Storage) Close() error {
	return s.storage.Close()
}

func (s *Storage) GetUserIDWithLogin(ctx context.Context, login string) (int, error) {
	query := `SELECT id FROM users WHERE login = $1`
	var userID int

	if err := s.storage.QueryRowContext(ctx, query, login).Scan(&userID); err != nil {
		s.log.Error("failed to check login in base", "error", err)
		return -1, err
	}

	return userID, nil
}

func (s *Storage) GetUserIDWithToken(ctx context.Context, token string) (int, error) {
	query := `SELECT id FROM users WHERE token = $1`
	var userID int

	if err := s.storage.QueryRowContext(ctx, query, token).Scan(&userID); err != nil {
		s.log.Error("failed to check login in base", "error", err)
		return -1, err
	}

	return userID, nil
}

// SaveLoginAndToken - сохраняет логин и токен в базе данных.
func (s *Storage) SaveLoginAndToken(ctx context.Context, login, token string) error {

	query := `INSERT INTO users (login, token) VALUE ($1, $2)`
	_, err := s.storage.ExecContext(ctx, query, login, token)
	if err != nil {
		s.log.Error("failed to update access token", "error", err)
		return err
	}
	s.log.Info("access token updated")
	return nil
}

// UpdateLoginAndToken - обновляет логин и токен в базе данных.
func (s *Storage) UpdateLoginAndToken(ctx context.Context, userID int, token string) error {
	now := time.Now()
	// Update the token for the existing user
	query := `UPDATE users SET token = $1, updated_at = $2 WHERE id = $3`
	_, err := s.storage.ExecContext(ctx, query, token, now, userID)
	if err != nil {
		s.log.Error("failed to update access token", "error", err)
		return err
	}
	s.log.Info("access token updated")
	return nil
}

// SaveLoginAndPasswordInCredentials - сохраняет полученный логин и пароль от ресурса в базу.
func (s *Storage) SaveLoginAndPasswordInCredentials(
	ctx context.Context,
	userID int,
	resource string,
	login string,
	password string) error {

	query := `INSERT INTO credentials (user_id, resource, login, password) VALUES ($1, $2, $3, $4)`

	_, err := s.storage.ExecContext(ctx, query, userID, resource, login, password)
	if err != nil {
		s.log.Error("failed to save in credentials", "error", err)
		return err
	}

	return nil
}

// SaveTextData - сохраняет получены текст в базу.
func (s *Storage) SaveTextDataInDatabase(ctx context.Context, data string) error {
	query := `INSERT INTO text_data (data) VALUES ($1)`

	_, err := s.storage.ExecContext(ctx, query, data)
	if err != nil {
		s.log.Error("failed to save in credentials", "error", err)
		return err
	}

	return nil
}

// SaveBinaryData - сохраняет полученные бинарные данные.
func (s *Storage) SaveBinaryDataInDatabase(ctx context.Context, data string) error {
	query := `INSERT INTO binary_data (data) VALUES ($1)`

	_, err := s.storage.ExecContext(ctx, query, data)
	if err != nil {
		s.log.Error("failed to save in credentials", "error", err)
		return err
	}

	return nil
}

// SaveCardsInDatabase - сохраняет полученные данные по картам в базу.
func (s *Storage) SaveCardsInDatabase(ctx context.Context, cards string) error {
	query := `INSERT INTO cards ( cards) VALUES ($1)`

	_, err := s.storage.ExecContext(ctx, query, cards)
	if err != nil {
		s.log.Error("failed to save in cards", "error", err)
		return err
	}

	return nil
}

// GetAll - возвращает все данные из базы данных.
func (s *Storage) GetAll(ctx context.Context, tableName string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := s.storage.QueryContext(ctx, query)
	if err != nil {
		s.log.Error("failed to get all data from database", "error", err)
		return nil, err
	}

	return rows, nil
}
