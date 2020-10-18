package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	var errorsCount int32
	var currentTaskNum int
	taskLength := len(tasks)
	for {
		wg := sync.WaitGroup{}
		wg.Add(N)
		for i := 0; i < N; i++ {
			if currentTaskNum >= taskLength {
				return nil
			}
			task := tasks[currentTaskNum]
			currentTaskNum += 1
			go func(task Task) {
				err := task()
				if err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
				wg.Done()
			}(task)
		}
		wg.Wait()
		if M <= 0 {
			continue
		}
		if errorsCount >= int32(M) {
			return ErrErrorsLimitExceeded
		}
	}
}
