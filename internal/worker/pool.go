package worker

import (
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

var ErrPoolClosed = errors.New("worker pool is closed")

const (
	ctxKeyWorkerCount       = ctxKey(KeyWorkerCount)
	ctxKeySubmittedJobs     = ctxKey(KeySubmittedJobs)
	ctxKeyFailedSubmissions = ctxKey(KeyFailedSubmissions)
	ctxKeyPoolStartedAt     = ctxKey(KeyPoolStartedAt)
	ctxKeyPoolStoppedAt     = ctxKey(KeyPoolStoppedAt)
	ctxKeyPoolDuration      = ctxKey(KeyPoolDuration)
	ctxKeyPoolClosed        = ctxKey(KeyPoolClosed)
)

const (
	KeyWorkerCount       = "worker_count"
	KeySubmittedJobs     = "jobs_submitted"
	KeyFailedSubmissions = "failed_submissions"
	KeyPoolStartedAt     = "pool_started_at"
	KeyPoolStoppedAt     = "pool_stopped_at"
	KeyPoolDuration      = "pool_duration_seconds"
	KeyPoolClosed        = "pool_closed"
)

// Pool represents a worker pool that processes jobs concurrently with a fixed number of workers.
type Pool struct {
	maxWorkers     int
	jobs           chan *Job
	results        chan *JobResult
	wg             *sync.WaitGroup
	closed         atomic.Bool
	quit           chan struct{}
	startedAt      time.Time
	stoppedAt      time.Time
	duration       time.Duration
	submitted      int
	failedToSubmit int
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
		quit:       make(chan struct{}),
	}
}

// Run initializes and starts all workers in the pool, ensuring they are ready to process jobs concurrently.
func (p *Pool) Run() {
	p.startedAt = time.Now()
	for i := 1; i <= p.maxWorkers; i++ {
		nw := NewWorker(i, p.jobs, p.results, p.quit)
		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done() // Signal completion when the goroutine exits
			w.Start()
		}(nw)
	}
}

// Submit adds a job to the pool's job queue.
func (p *Pool) Submit(job *Job) (err error) {
	job.SetSubmittedAt()
	if p.closed.Load() {
		return ErrPoolClosed
	}
	defer func() {
		if r := recover(); r != nil {
			err = ErrPoolClosed
			p.failedToSubmit++
			slog.With(slog.String(KeyJobID, job.ID)).Warn("Job queue closed, job not submitted")
		}
	}()
	p.jobs <- job
	p.submitted++
	return nil
}

// SubmitBatch submits a batch of jobs to the worker pool for processing.
func (p *Pool) SubmitBatch(jobs []*Job) (int, int, error) {
	submitted := 0
	failures := 0
	var errs error
	for _, job := range jobs {
		err := p.Submit(job)
		if err != nil {
			failures++
			slog.With(slog.String(KeyJobID, job.ID)).Warn("Job failed", slog.Any("error", err))
			errs = errors.Join(errs, err)
		} else {
			submitted++
		}
	}
	return submitted, failures, errs
}

// Shutdown gracefully stops the pool by closing the job queue, waiting for workers to complete,
// and closing the results channel.
func (p *Pool) Shutdown() {
	if p.closed.CompareAndSwap(false, true) {
		p.stoppedAt = time.Now()
		p.duration = p.stoppedAt.Sub(p.startedAt)
		close(p.jobs)
		p.wg.Wait()
		close(p.results)
	}
}

// Stop gracefully stops the worker pool by closing the job queue and waiting for all workers to finish their tasks.
// Allows for manual control over result retrieval.
func (p *Pool) Stop() {
	if p.closed.CompareAndSwap(false, true) {
		p.stoppedAt = time.Now()
		p.duration = p.stoppedAt.Sub(p.startedAt)
		close(p.jobs)
		p.wg.Wait()
	}
}

// Terminate signals the worker pool to terminate all workers immediately.
func (p *Pool) Terminate() {
	if p.closed.CompareAndSwap(false, true) {
		p.stoppedAt = time.Now()
		p.duration = p.stoppedAt.Sub(p.startedAt)
		// Cancel any ongoing work by closing channels immediately
		close(p.jobs)
		close(p.results)
	}
}

// Results returns a read-only channel to retrieve processed job results from the worker pool.
func (p *Pool) Results() <-chan *JobResult {
	return p.results
}

func (p *Pool) Duration() time.Duration {
	return p.duration
}

func (p *Pool) StartedAt() time.Time {
	return p.startedAt
}

func (p *Pool) StoppedAt() time.Time {
	return p.stoppedAt
}

func (p *Pool) LogValue() slog.Value {
	return slog.GroupValue(slog.Int(KeyWorkerCount, p.maxWorkers),
		slog.Int(KeySubmittedJobs, p.submitted),
		slog.Int(KeyFailedSubmissions, p.failedToSubmit),
		slog.Bool(KeyPoolClosed, p.closed.Load()),
		slog.Time(KeyPoolStartedAt, p.startedAt),
		slog.Time(KeyPoolStoppedAt, p.stoppedAt),
		slog.Float64(KeyPoolDuration, p.duration.Seconds()))
}
