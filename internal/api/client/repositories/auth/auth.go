package auth

import (
	"context"
	"google.golang.org/grpc"
	v1_pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

// service - интерфейс на сервисный слой
//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=auth
type service interface {
	SaveTokenInBase(ctx context.Context, login, password, token string) error
	CheckUser(ctx context.Context, login, password string) (string, error)
}

// SignalAuth - сигнал об авторизации
var SignalAuth = make(chan bool)

// Handlers - обработчик запросов
type Handlers struct {
	log     *slog.Logger
	service service
}

// NewHandlers - конструктор обработчика
func NewHandlers(log *slog.Logger, service service) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

// RegisterUser - регистрирует пользователя
func (h *Handlers) RegisterUser(ctx context.Context, conn *grpc.ClientConn, login, password string) error {
	// создаем клиента для регистрации
	registerClient := v1_pd.NewRegisterClient(conn)

	_, err := registerClient.Register(ctx, &v1_pd.RegisterRequest{
		Login:    login,
		Password: password,
	})

	if err != nil {
		h.log.Error("auth.repositories.app: failed to register user")
		return err
	}
	return nil
}

// AuthUser - авторизует пользователя
func (h *Handlers) AuthUser(ctx context.Context, conn *grpc.ClientConn, login, password string) (string, error) {
	// создаем клиента для авторизации
	authClient := v1_pd.NewAuthClient(conn)

	token, err := authClient.Auth(ctx, &v1_pd.AuthRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		h.log.Error("failed to auth user")
		return "", err
	}

	err = h.service.SaveTokenInBase(ctx, login, password, token.Token)
	if err != nil {
		h.log.Error("failed to handlers token", "error", err)
		return "", err
	}

	// сигнализируем об авторизации
	SignalAuth <- true

	return token.Token, nil
}

func (h *Handlers) AuthUserOffLine(ctx context.Context, login, password string) (string, error) {

	token, err := h.service.CheckUser(ctx, login, password)
	if err != nil {
		h.log.Error("failed to check user", "error", err)
		return "", err
	}
	return token, nil
}
