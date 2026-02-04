package download

import (
	"context"
	"sync"
)

// Job represents a download job
type Job struct {
	Bucket    string
	Key       string
	LocalPath string
	Size      int64
}

// Result represents a job result
type Result struct {
	Job   Job
	Error error
}

// WorkerPool manages a pool of download workers
type WorkerPool struct {
	workers int
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan Job, workers*2),
		results: make(chan Result, workers*2),
	}
}

// Start starts the worker pool
func (p *WorkerPool) Start(ctx context.Context, worker func(context.Context, Job) error) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					err := worker(ctx, job)
					select {
					case p.results <- Result{Job: job, Error: err}:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}
}

// Submit submits a job to the pool
func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job
}

// Results returns the results channel
func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

// Close closes the job channel and waits for workers to finish
func (p *WorkerPool) Close() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}

// Semaphore provides a simple semaphore implementation
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a new semaphore with the given capacity
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, n),
	}
}

// Acquire acquires a slot
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release releases a slot
func (s *Semaphore) Release() {
	<-s.ch
}

// TryAcquire tries to acquire a slot without blocking
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}
