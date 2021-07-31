package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(queue <-chan Task, errCnt *uint32, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range queue {
		err := task()
		if err != nil {
			atomic.AddUint32(errCnt, 1)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (err error) {
	var errCnt uint32
	limit := uint32(m)

	taskQueue := make(chan Task)
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(taskQueue, &errCnt, wg)
	}

	for _, task := range tasks {
		if limit > 0 && atomic.LoadUint32(&errCnt) >= limit {
			err = ErrErrorsLimitExceeded
			break
		}
		taskQueue <- task
	}

	close(taskQueue)
	wg.Wait()

	return err
}
