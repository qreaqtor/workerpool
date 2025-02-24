package main

import (
	"fmt"
	"sync"
	"time"

	workerpool "github.com/qreaqtor/workerpool/worker-pool"
)

const (
	workersCount = 5
	actionsCount = 10
)

func main() {
	wp := workerpool.NewWorkerPool(workersCount)
	wg := new(sync.WaitGroup)

	wg.Add(actionsCount)
	for i := range actionsCount {
		err := wp.Push(getPrintlnAction(wg, i + 1))
		if err != nil {
			fmt.Println(err)
		}
	}

	wg.Wait()
	err := wp.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func getPrintlnAction(wg *sync.WaitGroup, i int) workerpool.Action {
	return func() {
		time.Sleep(time.Duration(i * int(time.Second)))
		fmt.Println(i)
		wg.Done()
	}
}
