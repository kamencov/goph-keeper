package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

type service interface {
	ValidateToken(ctx context.Context, token string) (int, error)
	CreatTokenForUser(userID string) (string, error)
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

func (m *Middleware) UnaryInterceptor(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {

	// Логируем запрос с использованием вашего логгера
	start := time.Now()
	m.log.Info(
		"gRPC request started",
		"method", info.FullMethod,
		"request", req,
	)

	// Методы, которые требуют аутентификации
	methodsRequiringCheckAuth := map[string]bool{
		"/goph_keeper_v1.SyncFromClient/SyncFromClientCredentials": true,
		"/goph_keeper_v1.SyncFromClient/SyncFromClientTextData":    true,
		"/goph_keeper_v1.SyncFromClient/SyncFromClientBinaryData":  true,
		"/goph_keeper_v1.SyncFromClient/SyncFromClientCards":       true,
	}

	// Если метод в списке, применяем AuthCheckInterceptor
	if methodsRequiringCheckAuth[info.FullMethod] {
		var authErr error
		ctx, authErr = m.AuthCheckInterceptor(ctx) // Передаём обновлённый контекст
		if authErr != nil {
			m.log.Error(
				"authentication failed",
				"method", info.FullMethod,
				"error", authErr,
			)
			return nil, authErr
		}
	}

	// Вызываем основной обработчик
	resp, err = handler(ctx, req)
	duration := time.Since(start)

	// Логируем завершение обработки
	if err != nil {
		m.log.Error(
			"gRPC request failed",
			"method", info.FullMethod,
			"error", err,
			"duration", duration,
		)
		return nil, err
	}

	m.log.Info(
		"gRPC request completed",
		"method", info.FullMethod,
		"response", resp,
		"duration", duration,
	)

	return resp, nil
}

func (m *Middleware) AuthCheckInterceptor(ctx context.Context) (context.Context, error) {
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
