package workerpool

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Description string

	WorkersCount int
	ActionsCount int
	Sleep        time.Duration

	Expected time.Duration
}

func getSleepAction(wg *sync.WaitGroup, d time.Duration) Action {
	return func() {
		time.Sleep(d)
		wg.Done()
	}
}

const delta = time.Millisecond * 100

func runTestCase(wp *WorkerPool, testCase testCase, t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(testCase.ActionsCount)

	startAt := time.Now()
	for range testCase.ActionsCount {
		go func(sleep time.Duration) {
			err := wp.Push(getSleepAction(wg, sleep))
			assert.NoError(t, err)
		}(testCase.Sleep)
	}
	wg.Wait()
	actualTime := time.Since(startAt)

	assert.InDelta(t, testCase.Expected, actualTime, float64(delta))

	err := wp.Close()
	assert.NoError(t, err)
}

func TestWorkerPoolWithSleepAction(t *testing.T) {
	cpuCount := runtime.NumCPU()

	testCases := []testCase{
		{
			Description:  "Count of workers equal count of actions",
			WorkersCount: cpuCount,
			ActionsCount: cpuCount,
			Sleep:        time.Second,
			Expected:     time.Second,
		},
		{
			Description:  "One worker and many actions",
			WorkersCount: 1,
			ActionsCount: cpuCount,
			Sleep:        time.Second,
			Expected:     time.Duration(cpuCount * int(time.Second)),
		},
		{
			Description:  "Workers count 2 times less then actions ",
			WorkersCount: cpuCount,
			ActionsCount: cpuCount * 2,
			Sleep:        time.Second,
			Expected:     time.Second * 2,
		},
		{
			Description:  "Workers count should be number of CPU when input count is zero",
			WorkersCount: 0,
			ActionsCount: cpuCount,
			Sleep:        time.Second,
			Expected:     time.Second,
		},
		{
			Description:  "Workers count should be number of CPU when input count is negative number",
			WorkersCount: -10,
			ActionsCount: cpuCount,
			Sleep:        time.Second,
			Expected:     time.Second,
		},
	}

	for _, testCase := range testCases {
		wp := NewWorkerPool(testCase.WorkersCount)

		runTestCase(wp, testCase, t)
	}
}

func TestWorkerPoolPushAfterCloseShouldBeWithError(t *testing.T) {
	testCase := testCase{
		WorkersCount: 1,
		ActionsCount: 1,
		Sleep:        time.Second,
		Expected:     time.Second,
	}

	wp := NewWorkerPool(testCase.WorkersCount)

	runTestCase(wp, testCase, t)

	err := wp.Push(func() {})
	assert.ErrorIs(t, err, errPoolClosed)
}

func TestWorkerPoolAfterFirstCloseShouldBeWithError(t *testing.T) {
	testCase := testCase{
		WorkersCount: 1,
		ActionsCount: 1,
		Sleep:        time.Second,
		Expected:     time.Second,
	}

	wp := NewWorkerPool(testCase.WorkersCount)

	runTestCase(wp, testCase, t)

	err := wp.Close()
	assert.ErrorIs(t, err, errPoolClosed)
}
