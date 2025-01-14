package workers

import (
	"context"
	"google.golang.org/grpc"
	"goph-keeper/internal/api/client/repositories/auth"
	"sync"
	"time"
)

type service interface {
	PushData(ctx context.Context, conn *grpc.ClientConn) error
}

type Worker struct {
	service service
	conn    *grpc.ClientConn
	wg      *sync.WaitGroup
}

func NewWorker(service service, conn *grpc.ClientConn) *Worker {
	return &Worker{
		service: service,
		conn:    conn,
		wg:      &sync.WaitGroup{},
	}
}

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
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}
