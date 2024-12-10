package memory

import (
	"context"
	"log/slog"
	"sync"
)

type Memory struct {
	storage map[string]string
	mu      sync.RWMutex
	log     *slog.Logger
}

func NewMemory(log *slog.Logger) *Memory {
	return &Memory{
		storage: make(map[string]string),
		log:     log,
	}
}

// Close закрывает хранилище.
func (m *Memory) Close() error {
	return nil
}
func (m *Memory) CheckUser(ctx context.Context, login string) error {
	return nil
}
func (m *Memory) CheckPassword(login string) (string, bool) {
	return "", false
}
func (m *Memory) SaveUser(ctx context.Context, login, hashPassword string) error {
	return nil
}
func (m *Memory) SaveTableUserAndUpdateToken(login, accessToken string) error {
	return nil
}
func (m *Memory) SaveLoginAndPasswordInCredentials(ctx context.Context, userID int, resource string, login string, password string) error {
	return nil
}
func (m *Memory) SaveTextData(ctx context.Context, userID int, data string) error {
	return nil
}
func (m *Memory) SaveBinaryData(ctx context.Context, uid int, data string) error { return nil }
func (m *Memory) SaveCards(ctx context.Context, userID int, cards string) error {
	return nil
}
func (m *Memory) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	return -1, nil
}
func (m *Memory) GetUserIDByLogin(ctx context.Context, login string) (int, error) {
	return -1, nil
}
