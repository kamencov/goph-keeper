package workers

import (
	"context"
	"google.golang.org/grpc"
	"goph-keeper/internal/api/client/repositories/auth"
	"log/slog"
	"sync"
	"time"
)

// service - интерфейс сервиса
//
//go:generate mockgen -source=worker.go -destination=worker_mock.go -package=workers
type service interface {
	PushData(ctx context.Context, conn *grpc.ClientConn) error
}

// Worker - структура для работы с горутинами
type Worker struct {
	log     *slog.Logger
	service service
	conn    *grpc.ClientConn
	wg      *sync.WaitGroup
}

// NewWorker - конструктор
func NewWorker(log *slog.Logger, service service, conn *grpc.ClientConn) *Worker {
	return &Worker{
		log:     log,
		service: service,
		conn:    conn,
		wg:      &sync.WaitGroup{},
	}
}

// Run - горутина
func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case signalAuth := <-auth.SignalAuth:
			if signalAuth {
				w.wg.Add(1)
				go func() {
					defer w.wg.Done()
					w.Push(ctx)
				}()
			}
		case <-ctx.Done():
			w.wg.Wait()
			return
		}
	}
}

// Push - метод для отправки данных
func (w *Worker) Push(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.wg.Add(1)
			go func() {
				defer w.wg.Done()
				if err := w.service.PushData(ctx, w.conn); err != nil {
					// Обработка ошибки
					w.log.Error("failed to push data", "error", err)
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}
