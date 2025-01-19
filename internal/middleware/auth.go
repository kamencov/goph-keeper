package middleware

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

type service interface {
	ValidateToken(ctx context.Context, token string) (int, error)
}

type Middleware struct {
	log     *slog.Logger
	service service
}

func NewMiddleware(log *slog.Logger, service service) *Middleware {
	return &Middleware{
		log:     log,
		service: service,
	}
}

func (m *Middleware) AuthInterceptor(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	tokens := md["authorization"]
	if len(tokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	token := tokens[0]
	userID, err := m.service.ValidateToken(ctx, token) // Ваша логика проверки токена
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// Добавляем userID в контекст
	newCtx := context.WithValue(ctx, UserIDContextKey, userID)
	return newCtx, nil
}
