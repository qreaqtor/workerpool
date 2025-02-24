package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/qreaqtor/workerpool"
)

func getPrintlnAction(wg *sync.WaitGroup, i int) workerpool.Action {
	return func() {
		time.Sleep(time.Duration(i * int(time.Second)))
		fmt.Println(i)
		wg.Done()
	}
}

func main() {
	workersCount := 5
	actionsCount := 10

	wp := workerpool.NewWorkerPool(workersCount)
	wg := new(sync.WaitGroup)
	wg.Add(actionsCount)

	for i := range actionsCount {
		err := wp.Push(getPrintlnAction(wg, i))
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
