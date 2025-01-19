package postgresql

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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

	err = p.applyMigrations()
	if err != nil {
		return err
	}
	return nil
}

// applyMigrations - выполняет миграции через Goose
func (p *Postgresql) applyMigrations() error {
	migrationsDir := "./internal/migrations"
	if err := goose.Up(p.storage, migrationsDir); err != nil {
		p.log.Error("failed to apply migrations", "error", err)
		return err
	}
	p.log.Info("migrations applied successfully")
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
		p.log.Error("failed to handlers user", "error", err)
		return err
	}

	return nil
}

// SaveTableUserAndUpdateToken - сохраняет пользователя в базе данных.
func (p *Postgresql) SaveTableUserAndUpdateToken(login, accessToken string) error {
	query := "UPDATE users SET token = $1 WHERE login = $2"

	_, err := p.storage.Exec(query, accessToken, login)
	if err != nil {
		p.log.Error("failed to handlers access token", "error", err)
		return err
	}
	return nil
}

// ServerSaveLoginAndPasswordInCredentials - сохраняет полученный логин и пароль от ресурса в базу.
func (p *Postgresql) ServerSaveLoginAndPasswordInCredentials(
	ctx context.Context,
	userID int,
	resource string,
	login string,
	password string) error {
	query := `INSERT INTO credentials (user_id,resource, login, password) VALUES ($1, $2, $3, $4)`

	_, err := p.storage.ExecContext(ctx, query, userID, resource, login, password)
	if err != nil {
		p.log.Error("failed to handlers in credentials", "error", err)
		return err
	}

	return nil
}

// SaveTextDataPstgres - сохраняет получены текст в базу.
func (p *Postgresql) SaveTextDataPstgres(ctx context.Context, userID int, data string) error {
	query := `INSERT INTO text_data (user_id, text) VALUES ($1, $2)`

	_, err := p.storage.ExecContext(ctx, query, userID, data)
	if err != nil {
		p.log.Error("failed to handlers in credentials", "error", err)
		return err
	}

	return nil
}

// SaveBinaryData - сохраняет полученные бинарные данные.
func (p *Postgresql) SaveBinaryDataBinary(ctx context.Context, uid int, data string) error {
	query := `INSERT INTO binary_data (user_id, binary_data) VALUES ($1, $2)`

	_, err := p.storage.ExecContext(ctx, query, uid, data)
	if err != nil {
		p.log.Error("failed to handlers in credentials", "error", err)
		return err
	}

	return nil
}

// SaveCards - сохраняет полученные данные по картам в базу.
func (p *Postgresql) SaveCards(ctx context.Context, userID int, cards string) error {
	query := `INSERT INTO cards (user_id, cards) VALUES ($1, $2)`

	_, err := p.storage.ExecContext(ctx, query, userID, cards)
	if err != nil {
		p.log.Error("failed to handlers in cards", "error", err)
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

func (p *Postgresql) DeletedCredentials(ctx context.Context, userID int, resource string) error {
	query := `UPDATE credentials SET deleted = 1 WHERE user_id = $1 AND resource = $2`

	_, err := p.storage.ExecContext(ctx, query, userID, resource)
	if err != nil {
		p.log.Error("failed to deleted credentials", "error", err)
		return err
	}
	return nil
}

func (p *Postgresql) DeletedText(ctx context.Context, userID int, data string) error {
	query := `UPDATE text_data SET deleted = 1 WHERE user_id = $1 AND text = $2`

	_, err := p.storage.ExecContext(ctx, query, userID, data)
	if err != nil {
		p.log.Error("failed to deleted text", "error", err)
		return err
	}
	return nil
}

func (p *Postgresql) DeletedBinary(ctx context.Context, userID int, data string) error {
	query := `UPDATE binary_data SET deleted = 1 WHERE user_id = $1 AND binary_data = $2`

	_, err := p.storage.ExecContext(ctx, query, userID, data)
	if err != nil {
		p.log.Error("failed to deleted binary", "error", err)
		return err
	}
	return nil
}

func (p *Postgresql) DeletedCards(ctx context.Context, userID int, data string) error {
	query := `UPDATE cards SET deleted = 1 WHERE user_id = $1 AND cards = $2`

	_, err := p.storage.ExecContext(ctx, query, userID, data)
	if err != nil {
		p.log.Error("failed to deleted cards", "error", err)
		return err
	}
	return nil
}
