package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	var errorsCount int32
	var currentTaskNum int
	taskLength := len(tasks)
	for {
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := 0; i < n; i++ {
			if currentTaskNum >= taskLength {
				return nil
			}
			task := tasks[currentTaskNum]
			currentTaskNum++
			go func(task Task) {
				err := task()
				if err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
				wg.Done()
			}(task)
		}
		wg.Wait()
		if m <= 0 {
			continue
		}
		if errorsCount >= int32(m) {
			return ErrErrorsLimitExceeded
		}
	}
}
