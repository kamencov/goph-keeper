package credentials

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pd "goph-keeper/internal/proto/v1"
	"goph-keeper/internal/services/credentials"
	"log/slog"
)

type serviceCredentials interface {
	SaveLoginAndPassword(ctx context.Context, info, login, password string) error
}

type Handlers struct {
	pd.UnimplementedPostCredentialsServer
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

	err := h.service.SaveLoginAndPassword(ctx, in.GetResource(), in.GetLogin(), in.GetPassword())

	if err != nil {
		if errors.Is(err, credentials.ErrNotFoundUser) {
			h.log.Error("failed to get login in base", "error", err)
			return nil, status.Errorf(codes.NotFound, "login is not correct")
		}
		h.log.Error("failed to save login and password", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to save login and password")
	}

	return &pd.Empty{
		Message: "save complete",
	}, nil
}
