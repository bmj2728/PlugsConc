package worker

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"PlugsConc/internal/logger"
)

// ErrPoolClosed indicates that the worker pool has been closed and cannot accept any new jobs.
var ErrPoolClosed = errors.New("worker pool is closed")

// MetricResult represents the outcome of a job
type MetricResult struct {
	isSuccess bool
}
type BatchErrors map[string]error

func (b BatchErrors) Add(jobID string, err error) BatchErrors {
	b[jobID] = err
	return b
}

func (b BatchErrors) LogValue() slog.Value {
	var formatted strings.Builder
	if len(b) == 0 {
		return slog.AnyValue("")
	}
	formatted.WriteString("Batch Errors:\n")
	for j, e := range b {
		entry := fmt.Sprintf("%s: %s\n", j, e.Error())
		formatted.WriteString(entry)
	}
	batchGroup := slog.GroupValue(slog.String(logger.KeyBatchErrors, formatted.String()))
	return batchGroup
}

// NewMetricResult creates and returns a new MetricResult with the given success status.
func NewMetricResult(isSuccess bool) *MetricResult {
	return &MetricResult{
		isSuccess: isSuccess,
	}
}

// Pool represents a worker pool used to manage the execution of concurrent jobs.
type Pool struct {
	maxWorkers     int                // workers count
	jobs           chan *Job          // for incoming jobs
	results        chan *JobResult    // for completed jobs
	wg             *sync.WaitGroup    // for workers
	closed         atomic.Bool        // identify if closed
	quit           chan struct{}      // for quit signals
	metricsChannel chan *MetricResult // pool metrics chan
	metrics        *PoolMetrics       // pool metrics
}

// NewPool initializes a new Pool with the specified number of workers and a buffer size for its channels.
func NewPool(maxWorkers int, buffer int) *Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	var jobs chan *Job
	var results chan *JobResult
	var metricsConsumer chan *MetricResult
	if buffer < 1 {
		// create unbuffered channels
		jobs = make(chan *Job)
		results = make(chan *JobResult)
		metricsConsumer = make(chan *MetricResult)
	} else {
		// create buffered channels
		jobs = make(chan *Job, buffer)
		results = make(chan *JobResult, buffer)
		metricsConsumer = make(chan *MetricResult, buffer)
	}
	return &Pool{
		maxWorkers:     maxWorkers,
		jobs:           jobs,
		results:        results,
		wg:             &sync.WaitGroup{},
		quit:           make(chan struct{}),
		metricsChannel: metricsConsumer,
		metrics:        NewPoolMetrics(),
	}
}

// Run starts the worker pool and initializes the configured number of worker goroutines to process jobs concurrently.
func (p *Pool) Run() {
	p.metrics.SetStarted()
	go p.collectMetrics()
	for i := 1; i <= p.maxWorkers; i++ {
		nw := NewWorker(i, p.jobs, p.results, p.quit, p.metricsChannel)
		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done() // Signal completion when the goroutine exits
			w.Start()
		}(nw)
	}
}

// Submit schedules a Job for execution in the Pool; returns an error if the Pool is closed or the submission fails.
func (p *Pool) Submit(job *Job) (err error) {
	job.SetSubmittedAt()
	if p.closed.Load() {
		return ErrPoolClosed
	}
	defer func() {
		if r := recover(); r != nil {
			err = ErrPoolClosed
			p.metrics.RecordFailedSubmission()
			slog.With(slog.String(logger.KeyJobID, job.ID)).Warn("Job queue closed, job not submitted")
		}
	}()
	p.jobs <- job
	p.metrics.RecordSubmission()
	return nil
}

// SubmitBatch processes a batch of jobs, submitting each to the pool and tracking the number of successes and failures.
func (p *Pool) SubmitBatch(jobs []*Job) (int, int, BatchErrors) {
	submitted := 0
	failures := 0
	var errorMap BatchErrors
	for _, job := range jobs {
		err := p.Submit(job)
		if err != nil {
			failures++
			slog.With(slog.String(logger.KeyJobID, job.ID)).Warn("Job failed", slog.Any("error", err))
			errorMap.Add(job.ID, err)
		} else {
			submitted++
		}
	}
	return submitted, failures, errorMap
}

// Shutdown gracefully stops the worker pool, ensuring all submitted jobs are completed and resources are released.
func (p *Pool) Shutdown() {
	if p.closed.CompareAndSwap(false, true) {
		p.metrics.SetStopped()
		close(p.jobs)
		p.wg.Wait()
		p.metrics.SetCompleted()
		err := p.metrics.SetDuration()
		if err != nil {
			slog.Warn("unable to set metrics")
		}
		close(p.results)
		close(p.metricsChannel)
	}
}

// Stop gracefully shuts down the pool by marking it as closed, waiting for workers to finish, and finalizing metrics.
func (p *Pool) Stop() {
	if p.closed.CompareAndSwap(false, true) {
		p.metrics.SetStopped()
		close(p.jobs)
		p.wg.Wait()
		p.metrics.SetCompleted()
		err := p.metrics.SetDuration()
		if err != nil {
			slog.Warn("unable to set pool duration")
		}
	}
}

// Terminate gracefully stops the pool execution by closing all channels and setting metrics,
// canceling ongoing work immediately.
func (p *Pool) Terminate() {
	if p.closed.CompareAndSwap(false, true) {
		p.metrics.SetStopped()
		// Cancel any ongoing work by closing channels immediately
		close(p.jobs)
		p.metrics.SetCompleted()
		err := p.metrics.SetDuration()
		if err != nil {
			slog.Warn("unable to set pool duration")
		}
		close(p.results)
		close(p.metricsChannel)
	}
}

// Results returns a channel from which completed job results can be received.
func (p *Pool) Results() <-chan *JobResult {
	return p.results
}

// Duration returns the total duration for which the pool has been active, as tracked by its metrics.
func (p *Pool) Duration() time.Duration {
	return p.metrics.Duration()
}

// StartedAt returns the timestamp when the pool was started by accessing the started time from the pool's metrics.
func (p *Pool) StartedAt() time.Time {
	return p.metrics.Started()
}

// StoppedAt returns the timestamp when the pool was stopped, as recorded in the pool's metrics.
func (p *Pool) StoppedAt() time.Time {
	return p.metrics.Stopped()
}

func (p *Pool) CompletedAt() time.Time {
	return p.metrics.Completed()
}

// Workers returns the maximum number of workers configured for the pool.
func (p *Pool) Workers() int {
	return p.maxWorkers
}

// Metrics returns a copy of the current pool metrics, providing a snapshot of important runtime statistics.
func (p *Pool) Metrics() *PoolMetrics {
	// create a new struct to copy into
	mCopy := NewPoolMetrics()
	// lock the existing until complete
	p.metrics.mu.Lock()
	defer p.metrics.mu.Unlock()
	// copy data
	mCopy.startedAt = p.metrics.startedAt
	mCopy.stoppedAt = p.metrics.stoppedAt
	mCopy.completedAt = p.metrics.completedAt
	mCopy.duration = p.metrics.duration
	mCopy.submissions = p.metrics.submissions
	mCopy.submissionFailures = p.metrics.submissionFailures
	mCopy.succeeded = p.metrics.succeeded
	mCopy.failed = p.metrics.failed
	//return copy
	return mCopy
}

// LogValue generates a structured log representation of the pool's state, including its closed status,
// worker count, and metrics.
func (p *Pool) LogValue() slog.Value {
	return slog.GroupValue(slog.Bool(logger.KeyPoolClosed, p.closed.Load()),
		slog.Int(logger.KeyWorkerCount, p.maxWorkers),
		slog.Any(logger.KeyPoolMetrics, p.metrics.LogValue()),
	)
}

// collectMetrics processes metric results from the metricsChannel, updating success and failure counts
// in a thread-safe manner.
func (p *Pool) collectMetrics() {
	for mr := range p.metricsChannel {
		p.metrics.mu.Lock()
		if mr.isSuccess {
			p.metrics.succeeded++
		} else {
			p.metrics.failed++
		}
		p.metrics.mu.Unlock()
	}
}
