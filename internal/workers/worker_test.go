package workers

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"goph-keeper/internal/api/client/repositories/auth"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log := slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug}))
	db := NewMockservice(ctrl)
	worker := NewWorker(log, db, nil)

	if worker == nil {
		t.Errorf("Worker is nil")
	}
}

func TestWorker_Run(t *testing.T) {
	tests := []struct {
		name string
		sig  bool
	}{
		{
			name: "successful_run",
			sig:  true,
		},
		{
			name: "failed_run",
			sig:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := NewMockservice(ctrl)

			worker := NewWorker(log, mockService, nil)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go worker.Run(ctx)
			auth.SignalAuth <- tt.sig
		})
	}

}

func TestWorker_Push(t *testing.T) {
	tests := []struct {
		name        string
		seconds     int
		expectedErr error
	}{
		{
			name:    "successful_push",
			seconds: 2,
		},
		{
			name:        "failed_push",
			seconds:     2,
			expectedErr: errors.New("failed to push data"),
		},
		{
			name:    "timeout_push",
			seconds: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			log := slog.New(slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug}))

			mockService := NewMockservice(ctrl)
			worker := NewWorker(log, mockService, nil)

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(tt.seconds)*time.Second))
			defer cancel()

			// Ожидание вызова метода PushData
			mockService.EXPECT().PushData(ctx, gomock.Any()).Return(tt.expectedErr).AnyTimes()

			// Запускаем Push в отдельной горутине, так как метод блокирующий
			go worker.Push(ctx)

			// Ждём завершения контекста или тестового времени
			time.Sleep(3 * time.Second)
		})
	}
}
