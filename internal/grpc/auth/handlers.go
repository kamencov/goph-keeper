package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pd "goph-keeper/internal/proto/v1"
	"goph-keeper/internal/services/server/auth"
	"log/slog"
)

// serviceAuth - интерфейс на сервисный слой.
//
//go:generate mockgen -source=handlers.go -destination=handlers_mock.go -package=auth
type serviceAuth interface {
	Auth(login, password string) (auth.Tokens, error)
}

// Handlers - обработчик запросов.
type Handlers struct {
	pd.UnimplementedAuthServer
	log     *slog.Logger
	service serviceAuth
}

// NewHandlers создаёт новый Handlers.
func NewHandlers(log *slog.Logger, service serviceAuth) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

// Auth - авторизует пользователя.
// @Tags POST
// @Summary Авторизует пользователя.
// @Description Авторизует пользователя.
// @Accept json
// @Produce json
// @Param request body v1_pd.AuthRequest true "request"
// @Success 200 {object} v1_pd.AuthResponse
// @Failure 400 "password or login is empty"
// @Failure 404 "login is not correct"
// @Failure 401 "password is not correct"
// @Failure 500 "failed to auth user"
// @Router /goph_keeper_v1.Auth/Auth [post]
func (h *Handlers) Auth(ctx context.Context, in *pd.AuthRequest) (*pd.AuthResponse, error) {
	if in.Password == "" || in.Login == "" {
		h.log.Error("password or login is empty")
		return nil, status.Errorf(codes.InvalidArgument, "password or login is empty")
	}

	token, err := h.service.Auth(in.GetLogin(), in.GetPassword())

	if err != nil {
		if errors.Is(err, auth.ErrNotFoundLogin) {
			h.log.Error("failed to check login in base", "error", err)
			return nil, status.Errorf(codes.NotFound, "login is not correct")
		}
		if errors.Is(err, auth.ErrWrongPassword) {
			h.log.Error("failed do password match", "error", err)
			return nil, status.Errorf(codes.Unauthenticated, "password is not correct")
		}
		h.log.Error("failed to auth user", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to auth user")
	}

	if token.AccessToken == "" {
		h.log.Error("access token is empty")
		return nil, status.Errorf(codes.Internal, "access token is empty")
	}

	return &pd.AuthResponse{
		Token:   token.AccessToken,
		Message: "token created",
	}, nil

}
