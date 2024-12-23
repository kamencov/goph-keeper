package auth

import (
	"context"
	"google.golang.org/grpc"
	v1_pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

type service interface {
	SaveTokenInBase(ctx context.Context, login, token string) error
}

type Handlers struct {
	log     *slog.Logger
	service service
}

func NewHandlers(log *slog.Logger, service service) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

func (h *Handlers) RegisterUser(ctx context.Context, conn *grpc.ClientConn, login, password string) error {
	// создаем клиента для регистрации
	registerClient := v1_pd.NewRegisterClient(conn)

	_, err := registerClient.Register(ctx, &v1_pd.RegisterRequest{
		Login:    login,
		Password: password,
	})

	if err != nil {
		h.log.Error("auth.handlers.app: failed to register user")
		return err
	}
	return nil
}

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

	err = h.service.SaveTokenInBase(ctx, login, token.Token)
	if err != nil {
		h.log.Error("failed to save token", "error", err)
		return "", err
	}
	return token.Token, nil
}
