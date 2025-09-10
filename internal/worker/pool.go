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
	ctxKeyPoolCompletedAt   = ctxKey(KeyPoolCompletedAt)
	ctxKeyPoolDuration      = ctxKey(KeyPoolDuration)
	ctxKeyPoolClosed        = ctxKey(KeyPoolClosed)
	ctxKeySuccessfulJobs    = ctxKey(KeySuccessfulJobs)
	ctxKeyFailedJobs        = ctxKey(KeyFailedJobs)
	ctxKeySPoolMetrics      = ctxKey(KeyPoolMetrics)
)

const (
	KeyWorkerCount       = "worker_count"
	KeySubmittedJobs     = "jobs_submitted"
	KeyFailedSubmissions = "failed_submissions"
	KeyPoolStartedAt     = "pool_started_at"
	KeyPoolStoppedAt     = "pool_stopped_at"
	KeyPoolCompletedAt   = "pool_completed_at"
	KeyPoolDuration      = "pool_duration_seconds"
	KeyPoolClosed        = "pool_closed"
	KeySuccessfulJobs    = "successful_jobs"
	KeyFailedJobs        = "failed_jobs"
	KeyPoolMetrics       = "pool_metrics"
)

type MetricResult struct {
	isSuccess bool
}

func NewMetricResult(isSuccess bool) *MetricResult {
	return &MetricResult{
		isSuccess: isSuccess,
	}
}

// Pool represents a worker pool that processes jobs concurrently with a fixed number of workers.
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

// NewPool creates and initializes a new Pool with the specified number of workers
// and an optional buffer size for channels.
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

// Run initializes and starts all workers in the pool, ensuring they are ready to process jobs concurrently.
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

// Submit adds a job to the pool's job queue.
func (p *Pool) Submit(job *Job) (err error) {
	job.SetSubmittedAt()
	if p.closed.Load() {
		return ErrPoolClosed
	}
	defer func() {
		if r := recover(); r != nil {
			err = ErrPoolClosed
			p.metrics.RecordFailedSubmission()
			slog.With(slog.String(KeyJobID, job.ID)).Warn("Job queue closed, job not submitted")
		}
	}()
	p.jobs <- job
	p.metrics.RecordSubmission()
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

// Stop gracefully stops the worker pool by closing the job queue and waiting for all workers to finish their tasks.
// Allows for manual control over result retrieval.
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

// Terminate signals the worker pool to terminate all workers immediately.
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

// Results returns a read-only channel to retrieve processed job results from the worker pool.
func (p *Pool) Results() <-chan *JobResult {
	return p.results
}

func (p *Pool) Duration() time.Duration {
	return p.metrics.Duration()
}

func (p *Pool) StartedAt() time.Time {
	return p.metrics.Started()
}

func (p *Pool) StoppedAt() time.Time {
	return p.metrics.Stopped()
}

func (p *Pool) CompletedAt() time.Time {
	return p.metrics.Completed()
}

func (p *Pool) LogValue() slog.Value {
	return slog.GroupValue(slog.Bool(KeyPoolClosed, p.closed.Load()),
		slog.Int(KeyWorkerCount, p.maxWorkers),
		slog.Any(KeyPoolMetrics, p.metrics.LogValue()),
	)
}

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
