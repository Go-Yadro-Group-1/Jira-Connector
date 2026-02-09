package sync

import (
	"context"
	"fmt"
	"sync"
)

type TaskFunc[T any] func(ctx context.Context) (T, error)

type Result[T any] struct {
	Value T
	Err   error
	ID    string
}

type WorkerPool[T any] struct {
	maxWorkers int
	tasks      chan TaskFunc[T]
	results    chan Result[T]
	wg         sync.WaitGroup
}

func New[T any](maxWorkers int) *WorkerPool[T] {
	return &WorkerPool[T]{
		maxWorkers: maxWorkers,
		tasks:      make(chan TaskFunc[T], maxWorkers),
		results:    make(chan Result[T], maxWorkers),
	}
}

func (p *WorkerPool[T]) Start(ctx context.Context) {
	p.wg.Add(p.maxWorkers)

	for i := 0; i < p.maxWorkers; i++ {
		go p.worker(ctx, i)
	}
}

func (p *WorkerPool[T]) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}

			result := Result[T]{ID: fmt.Sprintf("worker-%d", id)}
			result.Value, result.Err = task(ctx)

			p.results <- result
		}
	}
}

func (p *WorkerPool[T]) Submit(task TaskFunc[T]) {
	p.tasks <- task
}

func (p *WorkerPool[T]) Results() <-chan Result[T] {
	return p.results
}

func (p *WorkerPool[T]) Stop() {
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}
