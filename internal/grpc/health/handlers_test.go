package health

import (
	"context"
	v1_pd "goph-keeper/internal/proto/v1"
	"log/slog"
	"os"
	"testing"
)

func TestNewHandler(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	handler := NewHandler(log)

	if handler == nil {
		t.Errorf("Handler is nil")
	}
}

func TestHandler_Health(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

	handler := NewHandler(log)

	_, err := handler.Health(context.Background(), &v1_pd.Empty{})
	if err != nil {
		t.Errorf("Health() error = %v", err)
	}
}
