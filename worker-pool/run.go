package workerpool

import "context"

func (w *WorkerPool) run(ctx context.Context, workersCount int) {
	w.wg.Add(workersCount)

	for range workersCount {
		go w.runWorker(ctx)
	}
}

func (w *WorkerPool) runWorker(ctx context.Context) {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case action, ok := <-w.input:
			if !ok {
				return
			}
			action()
		}
	}
}
