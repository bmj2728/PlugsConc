package worker

import (
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
)

var ErrPoolClosed = errors.New("worker pool is closed")

// Pool represents a worker pool that processes jobs concurrently with a fixed number of workers.
type Pool struct {
	maxWorkers int
	jobs       chan *Job
	results    chan *JobResult
	wg         *sync.WaitGroup
	closed     atomic.Bool
}

// NewPool creates and initializes a new Pool with the specified number of workers
// and an optional buffer size for channels.
func NewPool(maxWorkers int, buffer int) *Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	var jobs chan *Job
	var results chan *JobResult
	if buffer < 1 {
		// create unbuffered channels
		jobs = make(chan *Job)
		results = make(chan *JobResult)
	} else {
		// create buffered channels
		jobs = make(chan *Job, buffer)
		results = make(chan *JobResult, buffer)
	}
	return &Pool{
		maxWorkers: maxWorkers,
		jobs:       jobs,
		results:    results,
		wg:         &sync.WaitGroup{},
	}
}

// Run initializes and starts all workers in the pool, ensuring they are ready to process jobs concurrently.
func (p *Pool) Run() {
	for i := 1; i <= p.maxWorkers; i++ {
		nw := NewWorker(i, p.jobs, p.results)
		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done() // Signal completion when the goroutine exits
			nw.Start()
		}(nw)
	}
}

// Submit adds a job to the pool's job queue.
func (p *Pool) Submit(job *Job) (err error) {
	if p.closed.Load() {
		return ErrPoolClosed
	}
	defer func() {
		if r := recover(); r != nil {
			err = ErrPoolClosed
			slog.With(slog.String(JobIDKey, job.ID)).Warn("Job queue closed, job not submitted")
		}
	}()
	p.jobs <- job
	return nil
}

// SubmitBatch submits a batch of jobs to the worker pool for processing.
func (p *Pool) SubmitBatch(jobs []*Job) (int, int, error) {
	processed := 0
	failures := 0
	var errs error
	for _, job := range jobs {
		err := p.Submit(job)
		if err != nil {
			failures++
			errs = errors.Join(errs, err)
		} else {
			processed++
		}
	}
	return processed, failures, errs
}

// Shutdown gracefully stops the pool by closing the job queue, waiting for workers to complete,
// and closing the results channel.
func (p *Pool) Shutdown() {
	if p.closed.CompareAndSwap(false, true) {
		close(p.jobs)
		p.wg.Wait()
		close(p.results)
	}
}

// Stop gracefully stops the worker pool by closing the job queue and waiting for all workers to finish their tasks.
// Allows for manual control over result retrieval.
func (p *Pool) Stop() {
	if p.closed.CompareAndSwap(false, true) {
		close(p.jobs)
		p.wg.Wait()
		close(p.results)
	}
}

// Terminate closes the job and result channels, ceasing all job submissions and result retrievals.
func (p *Pool) Terminate() {
	if p.closed.CompareAndSwap(false, true) {
		close(p.jobs)
		close(p.results)
	}
}

// Results returns a read-only channel to retrieve processed job results from the worker pool.
func (p *Pool) Results() <-chan *JobResult {
	return p.results
}
