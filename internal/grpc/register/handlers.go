package register

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pd "goph-keeper/internal/proto/v1"
	"goph-keeper/internal/storage/postgresql"
	"log/slog"
)

// serviceRegister - интерфейс на сервисный слой.
type serviceRegister interface {
	RegisterUser(ctx context.Context, login, password string) (int, error)
}

// Handlers структура в которой используем gRPC и сервисный слой.
type Handlers struct {
	pd.UnimplementedRegisterServer
	log     *slog.Logger
	service serviceRegister
}

// NewHandlers - конструктор регистрации пользователя.
func NewHandlers(log *slog.Logger, service serviceRegister) *Handlers {
	return &Handlers{
		log:     log,
		service: service}
}

// Register - регистрация пользователя.
func (h *Handlers) Register(ctx context.Context, in *pd.RegisterRequest) (*pd.RegisterResponse, error) {
	if in.Login == "" || in.Password == "" {
		h.log.Error("password or login is empty")
		return nil, status.Errorf(codes.InvalidArgument, "password or login is empty")
	}

	uid, err := h.service.RegisterUser(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		if errors.Is(err, postgresql.ErrUserAlreadyExists) {
			h.log.Error("failed, the user already exists", "error", err)
			return nil, status.Errorf(codes.AlreadyExists, "the user already exists")
		}
		h.log.Error("failed to register user", err)
		return nil, status.Errorf(codes.Internal, "failed to register user")
	}

	return &pd.RegisterResponse{
		Uid:     int32(uid),
		Message: "register completed"}, nil
}
