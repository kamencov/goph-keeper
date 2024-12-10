package postgresql

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"os"
)

var (
	ErrUserAlreadyExists = errors.New("the user already exists")
)

// Postgresql - подключение к базе данных.
type Postgresql struct {
	storage *sql.DB
	log     *slog.Logger
}

// NewPostgresql - конструктор, который возвращает Postgresql и error.
func NewPostgresql(log *slog.Logger) (*Postgresql, error) {
	p := &Postgresql{
		log: log,
	}
	err := p.initDB()
	return p, err
}

// initDB - инициализация подключения к базе данных и возвращает error.
func (p *Postgresql) initDB() error {
	path := os.Getenv("DATABASE_DSN")
	db, err := sql.Open("pgx", path)
	if err != nil {
		p.log.Error("failed to connect to database", "error", err)
		return err
	}

	p.storage = db
	p.log.Info("connected to database")

	err = p.createTableIfNotExists()
	if err != nil {
		return err
	}
	return nil
}

// createTableIfNotExists - создает таблицы в базе данных, если они не существуют.
func (p *Postgresql) createTableIfNotExists() (err error) {
	// Открытие транзакции
	tx, err := p.storage.Begin()
	if err != nil {
		p.log.Error("failed to begin transaction", "error", err)
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
        id SERIAL PRIMARY KEY, 
        login VARCHAR(255) NOT NULL, 
        password VARCHAR(255) NOT NULL, 
        token TEXT
    )`
	_, err = tx.Exec(query)
	if err != nil {
		p.log.Error("failed to create table - users", "error", err)
		return err
	}

	// Создаем таблицу credentials
	query = `CREATE TABLE IF NOT EXISTS credentials (
        id SERIAL PRIMARY KEY, 
        user_id INT NOT NULL, 
        resource VARCHAR(255) NOT NULL,
    	login VARCHAR(255) NOT NULL,
    	password VARCHAR(255) NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`
	_, err = tx.Exec(query)
	if err != nil {
		p.log.Error("failed to create table - credentials", "error", err)
		return err
	}

	// Создаем таблицу text_data
	query = `CREATE TABLE IF NOT EXISTS text_data (
        id SERIAL PRIMARY KEY, 
        user_id INT NOT NULL, 
        text TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`
	_, err = tx.Exec(query)
	if err != nil {
		p.log.Error("failed to create table - text_data", "error", err)
		return err
	}

	// Создаем таблицу binary_data
	query = `CREATE TABLE IF NOT EXISTS binary_data (
        id SERIAL PRIMARY KEY, 
        user_id INT NOT NULL, 
        binary BYTEA NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`
	_, err = tx.Exec(query)
	if err != nil {
		p.log.Error("failed to create table - binary_data", "error", err)
		return err
	}

	// Создаем таблицу cards
	query = `CREATE TABLE IF NOT EXISTS cards (
        id SERIAL PRIMARY KEY, 
        user_id INT NOT NULL, 
        cards TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`
	_, err = tx.Exec(query)
	if err != nil {
		p.log.Error("failed to create table - cards", "error", err)
		return err
	}

	// Подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		p.log.Error("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

// Close закрывает соединение с базой данных.
func (p *Postgresql) Close() error {
	return p.storage.Close()
}

// CheckUser - проверяет, существует ли пользователь в базе данных.
func (p *Postgresql) CheckUser(ctx context.Context, login string) error {
	query := "SELECT login FROM users WHERE login = $1"

	var foundLogin string
	err := p.storage.QueryRowContext(ctx, query, login).Scan(&foundLogin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Пользователь не найден, это не ошибка
			return nil
		}
		// Логируем ошибку при выполнении запроса
		p.log.Error("failed to check user", "error", err)
		return err
	}

	// Пользователь найден, возвращаем ошибку существующего пользователя
	return ErrUserAlreadyExists

}

// CheckPassword - проверяет, существует ли пользователь в базе данных.
func (p *Postgresql) CheckPassword(login string) (string, bool) {
	query := "SELECT password FROM users WHERE login = $1"

	row := p.storage.QueryRow(query, login)
	if row.Err() != nil {
		p.log.Error("failed to check password", "error", row.Err())
		return "", false
	}

	var passwordHash string
	if err := row.Scan(&passwordHash); err != nil {
		p.log.Error("failed to scan password", "error", err)
		return "", false
	}
	return passwordHash, true
}

// SaveUser - сохраняет пользователя в базе данных.
func (p *Postgresql) SaveUser(ctx context.Context, login, hashPassword string) error {
	query := "INSERT INTO users (login, password) VALUES ($1, $2)"
	_, err := p.storage.Exec(query, login, hashPassword)
	if err != nil {
		p.log.Error("failed to save user", "error", err)
		return err
	}

	return nil
}

// SaveTableUserAndUpdateToken - сохраняет пользователя в базе данных.
func (p *Postgresql) SaveTableUserAndUpdateToken(login, accessToken string) error {
	query := "UPDATE users SET token = $1 WHERE login = $2"

	_, err := p.storage.Exec(query, accessToken, login)
	if err != nil {
		p.log.Error("failed to save access token", "error", err)
		return err
	}
	return nil
}

// SaveLoginAndPasswordInCredentials - сохраняет полученный логин и пароль от ресурса в базу.
func (p *Postgresql) SaveLoginAndPasswordInCredentials(
	ctx context.Context,
	userID int,
	resource string,
	login string,
	password string) error {
	query := `INSERT INTO credentials (user_id,resource, login, password) VALUES ($1, $2, $3, $4)`

	_, err := p.storage.ExecContext(ctx, query, userID, resource, login, password)
	if err != nil {
		p.log.Error("failed to save in credentials", "error", err)
		return err
	}

	return nil
}

// SaveTextData - сохраняет получены текст в базу.
func (p *Postgresql) SaveTextData(ctx context.Context, userID int, data string) error {
	query := `INSERT INTO text_data (user_id, data) VALUES ($1, $2)`

	_, err := p.storage.ExecContext(ctx, query, userID, data)
	if err != nil {
		p.log.Error("failed to save in credentials", "error", err)
		return err
	}

	return nil
}

// SaveBinaryData - сохраняет полученные бинарные данные.
func (p *Postgresql) SaveBinaryData(ctx context.Context, uid int, data string) error {
	query := `INSERT INTO binary_data (user_id, data) VALUES ($1, $2)`

	_, err := p.storage.ExecContext(ctx, query, uid, data)
	if err != nil {
		p.log.Error("failed to save in credentials", "error", err)
		return err
	}

	return nil
}

// SaveCards - сохраняет полученные данные по картам в базу.
func (p *Postgresql) SaveCards(ctx context.Context, userID int, cards string) error {
	query := `INSERT INTO cards (user_id, cards) VALUES ($1, $2)`

	_, err := p.storage.ExecContext(ctx, query, userID, cards)
	if err != nil {
		p.log.Error("failed to save in cards", "error", err)
		return err
	}

	return nil
}

// GetUserIDByToken - получает user_id по токену.
func (p *Postgresql) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	query := "SELECT id FROM users WHERE token = $1"

	var userID int
	err := p.storage.QueryRowContext(ctx, query, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Пользователь не найден, это не ошибка
			return -1, nil
		}
		// Логируем ошибку при выполнении запроса
		p.log.Error("failed to check user", "error", err)
		return -1, err
	}

	return userID, nil
}

// GetUserIDByLogin - получает user_id по логину.
func (p *Postgresql) GetUserIDByLogin(ctx context.Context, login string) (int, error) {
	query := `SELECT id FROM users WHERE login = $1`

	var uid int

	err := p.storage.QueryRowContext(ctx, query, login).Scan(&uid)

	if err := err; err != nil {
		p.log.Error("failed to get id", "error", err)
		return -1, sql.ErrNoRows
	}

	return uid, nil
}
