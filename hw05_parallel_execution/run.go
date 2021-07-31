package hw05parallelexecution

import (
	"errors"
	"math"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	wg.Add(n)

	if m < 0 {
		m = 0
	}

	errChan := make(chan error, m)
	defer close(errChan)

	chunkSize := int(math.Ceil(float64(len(tasks)) / float64(n)))
	for i := 0; i < n; i++ {
		tasksChunk := getTasksChunk(i, chunkSize, tasks)
		go func() {
			defer wg.Done()
		loop:
			for _, task := range tasksChunk {
				err := task()
				if err != nil {
					select {
					case errChan <- err:
					default:
						break loop
					}
				}
			}
		}()
	}

	wg.Wait()

	select {
	case errChan <- errors.New(""):
	default:
		return ErrErrorsLimitExceeded
	}

	return nil
}

func getTasksChunk(i int, chunkSize int, tasks []Task) []Task {
	pos := i * chunkSize
	length := len(tasks)
	if pos >= length {
		return []Task{}
	}
	if pos+chunkSize > length {
		chunkSize = length - pos - 1
	}
	return tasks[pos : pos+chunkSize]
}
