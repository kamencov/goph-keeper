package workers

import (
	"context"
	"sync"
	"time"
)

type service interface {
	PushData()
}

type Worker struct {
	service service
	wg      *sync.WaitGroup
}

func NewWorker(service service) *Worker {
	return &Worker{
		service: service,
		wg:      &sync.WaitGroup{},
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.wg.Add(1)
			go func() {
				defer w.wg.Done()
				w.service.PushData()
			}()
		case <-ctx.Done():
			w.wg.Wait()
			return
		}
	}
}
