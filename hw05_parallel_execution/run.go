package hw05parallelexecution

import (
	"errors"
	"sync"
)

type errorState struct {
	mu           sync.Mutex
	errorCount   int
	errorLimit   int
	errorChannel chan struct{}
}

func newRunState(maxErrors int) *errorState {
	return &errorState{
		errorLimit:   maxErrors,
		errorChannel: make(chan struct{}),
	}
}

func (r *errorState) incErrorCnt() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.errorCount++
	if r.errorCount == r.errorLimit {
		close(r.errorChannel)
	}
}

type Task func() error

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	wg                     sync.WaitGroup
	state                  *errorState
)

func producer(tasks *[]Task, n int) chan Task {
	tasksChannel := make(chan Task, n)

	go func() {
		defer wg.Done()
		defer close(tasksChannel)

		for i := 0; i < len(*tasks); {
			select {
			case <-state.errorChannel:
				return
			default:
				if cap(tasksChannel) > len(tasksChannel) {
					tasksChannel <- (*tasks)[i]
					i++
				}
			}
		}
	}()

	return tasksChannel
}

func consumer(tasksChannel <-chan Task) {
	defer wg.Done()
	for {
		select {
		case <-state.errorChannel:
			return
		default:
			task, ok := <-tasksChannel
			if !ok {
				return
			}

			if err := task(); err != nil {
				state.incErrorCnt()
			}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	state = newRunState(m)

	// producer
	tasksChannel := producer(&tasks, n)

	// consumers
	for i := 0; i < n; i++ {
		go consumer(tasksChannel)
	}

	wg.Add(n + 1)
	wg.Wait()

	if state.errorLimit > 0 && state.errorCount >= state.errorLimit {
		return ErrErrorsLimitExceeded
	}
	return nil
}
