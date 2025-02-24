package workerpool

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
)

type WorkerPool struct {
	cancel context.CancelFunc

	input chan Action

	wg *sync.WaitGroup

	closed atomic.Bool
}

func NewWorkerPool(workersCount int) *WorkerPool {
	if workersCount <= 0 {
		workersCount = runtime.NumCPU()
	}

	wg := new(sync.WaitGroup)

	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		wg:     wg,
		cancel: cancel,
		input:  make(chan Action, workersCount),
	}

	wp.run(ctx, workersCount)

	return wp
}
