package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func NewWorkerPool(goroutines int, taskErrors int) *WorkerPool {
	pool := &WorkerPool{
		taskCh:     make(chan Task),
		doneCh:     make(chan struct{}),
		completeCh: make(chan struct{}),
		errCh:      make(chan error),
		goroutines: goroutines,
		taskErrors: taskErrors,
	}
	pool.startWorkers()
	return pool
}

type WorkerPool struct {
	taskCh       chan Task
	doneCh       chan struct{}
	completeCh   chan struct{}
	errCh        chan error
	ExecutionErr error
	goroutines   int
	taskErrors   int
	errorsCount  int32
	sync.RWMutex
	ready bool
}

func (e *WorkerPool) startWorkers() {
	wg := &sync.WaitGroup{}
	wg.Add(e.goroutines)
	go e.runErrorHandler()
	go e.runWorkersHandler(wg)
	for i := 0; i < e.goroutines; i++ {
		e.startWorker(wg, i)
	}
	e.ready = true
}
func (e *WorkerPool) runErrorHandler() {
	for {
		select {
		case <-e.errCh:
			if e.taskErrors <= 0 {
				continue
			}
			if atomic.AddInt32(&e.errorsCount, 1) >= int32(e.taskErrors) {
				e.stop(ErrErrorsLimitExceeded)
			}
		case <-e.completeCh:
			return
		}
	}
}

func (e *WorkerPool) runWorkersHandler(wg *sync.WaitGroup) {
	wg.Wait()
	fmt.Println("goroutines stopped")
	close(e.completeCh)
}

func (e *WorkerPool) startWorker(wg *sync.WaitGroup, i int) {
	go func(i int) {
		fmt.Println("starting worker:", i)
		for {
			select {
			case <-e.doneCh:
				fmt.Println("stop worker:", i)
				wg.Done()
				return
			case task := <-e.taskCh:
				if err := task(); err != nil {
					e.errCh <- err
				}
			}
		}
	}(i)
}
func (e *WorkerPool) stop(err error) {
	e.Lock()
	defer e.Unlock()
	if err != nil {
		e.ExecutionErr = err
	}
	if e.ready {
		close(e.doneCh)
		e.ready = false
	}
}

func (e *WorkerPool) ExecuteTask(task Task) bool {
	e.RLock()
	ready := e.ready
	e.RUnlock()

	if !ready {
		return false
	}
	select {
	case e.taskCh <- task:
		return true
	case <-e.doneCh:
	}
	return false
}

func (e *WorkerPool) Complete() error {
	e.stop(nil)
	<-e.completeCh
	return e.getError()
}

func (e *WorkerPool) getError() error {
	e.RLock()
	defer e.RUnlock()
	return e.ExecutionErr
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) (err error) {
	pool := NewWorkerPool(n, m)
	for _, task := range tasks {
		if !pool.ExecuteTask(task) {
			break
		}
	}
	return pool.Complete()
}
