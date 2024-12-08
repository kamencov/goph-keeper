package credentials

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pd "goph-keeper/internal/proto/v1"
	"log/slog"
)

type serviceCredentials interface {
	SaveLoginAndPassword(info, login, password string) error
}

type Handlers struct {
	pd.UnimplementedPostServer
	log     *slog.Logger
	service serviceCredentials
}

func NewHandlers(log *slog.Logger, service serviceCredentials) *Handlers {
	return &Handlers{
		log:     log,
		service: service,
	}
}

// PostLoginAndPassword сохраняет логин и пароль
func (h *Handlers) PostLoginAndPassword(ctx context.Context, in *pd.PostLoginAndPasswordRequest) (*pd.Empty, error) {

	if in.Password == "" || in.Login == "" {
		h.log.Error("password or login is empty")
		return nil, status.Errorf(codes.InvalidArgument, "password or login is empty")
	}

	err := h.service.SaveLoginAndPassword(in.Resource, in.Login, in.Password)

	if err != nil {
		h.log.Error("failed to save login and password", err)
		return nil, status.Errorf(codes.Internal, "failed to save login and password")
	}

	return &pd.Empty{
		Message: "save complete",
	}, nil
}
