package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pd "goph-keeper/internal/proto/v1"
	"goph-keeper/internal/services/auth"
	"log/slog"
)

type serviceAuth interface {
	Auth(login, password string) (auth.Tokens, error)
}

type Handlers struct {
	pd.UnimplementedAuthServer
	log     *slog.Logger
	service serviceAuth
}

func NewHandlers(log *slog.Logger, service serviceAuth) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

func (h *Handlers) Auth(ctx context.Context, in *pd.AuthRequest) (*pd.AuthResponse, error) {
	if in.Password == "" || in.Login == "" {
		h.log.Error("password or login is empty")
		return nil, status.Errorf(codes.InvalidArgument, "password or login is empty")
	}

	token, err := h.service.Auth(in.GetLogin(), in.GetPassword())

	if err != nil {
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
