package worker

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/bmj2728/PlugsConc/internal/logger"
	"github.com/hashicorp/go-hclog"
)

// Worker represents a worker that processes jobs from the jobs channel and sends results
// to the results channel.
type Worker struct {
	workerLogger hclog.Logger
	id           int
	jobs         <-chan *Job
	results      chan<- *JobResult
	metrics      chan<- *MetricResult
	quit         chan struct{}
}

// NewWorker creates and initializes a new Worker with a unique ID, a channel of jobs to process,
// and a results channel.
func NewWorker(id int, jobs <-chan *Job,
	results chan<- *JobResult,
	quit chan struct{},
	metrics chan<- *MetricResult,
	workerLogger hclog.Logger) *Worker {
	if workerLogger == nil {
		workerLogger = hclog.Default()
	}
	return &Worker{
		workerLogger: workerLogger,
		id:           id,
		jobs:         jobs,
		results:      results,
		quit:         quit,
		metrics:      metrics,
	}
}

// Start begins the worker's execution loop, processing jobs from the channel and sending results
// to the results channel.
func (w *Worker) Start() {
	w.workerLogger.Debug("Worker started")
	defer w.workerLogger.Debug("Worker stopped")

	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				return
			}
			// annotate job context
			job.Ctx = WithWorkerID(job.Ctx, w.id)
			job.SetStartedAt()

			// ensure cancellation and panic safety
			resultVal, err := func() (val any, err error) {
				// choose which cancel func to call on exit
				if job.CancelWithCause != nil {
					// capture the final err as the cause
					defer func() { job.CancelWithCause(err) }()
				} else if job.Cancel != nil {
					defer job.Cancel()
				}

				// panic safety: convert panics to errors
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("panic: %v\nstack: %s", r, string(debug.Stack()))
					}
				}()

				// retry loop
				delay := time.Duration(job.RetryDelay) * time.Millisecond
				for attempts := 0; ; attempts++ {
					job.Metrics.Attempts = attempts

					// if the job context is canceled, return immediately
					//  the default case is to continue the loop
					select {
					case <-job.Ctx.Done():
						job.SetFinishedAt()
						return nil, job.Ctx.Err()
					default:
					}

					// execute the job
					v, e := job.Execute(job.Ctx)
					// if the job succeeded, or we've reached the max retries, return the result/error
					//  otherwise, retry the job with a delay between retries'
					if e == nil || attempts >= job.MaxRetries {
						job.SetFinishedAt()
						return v, e
					}

					// log retry
					w.workerLogger.
						With(logger.KeyJobID, job.ID).
						With(logger.KeyRetryCount, attempts+1).
						Warn("Retrying job")

					// wait for the retry delay before continuing the loop
					if delay > 0 {
						t := time.NewTimer(delay)
						// if the job context is canceled, stop the timer and return immediately,
						//  otherwise, wait for the timer to expire
						select {
						case <-job.Ctx.Done():
							t.Stop()
							job.SetFinishedAt()
							return nil, job.Ctx.Err()
						case <-t.C:
						}
					}
				}
			}()

			// Safely send the result or quit if the pool is terminated.
			select {
			case w.results <- NewJobResult(job, w.id, resultVal, err):
				w.metrics <- NewMetricResult(err == nil)
				// Result sent successfully.
			case <-w.quit:
				// Pool was terminated while trying to send the result.
				// Log that the result is being discarded and exit the worker.
				job.SetFinishedAt()
				w.workerLogger.Warn("Worker terminated before sending result")
				return
			}

			attrs := []any{logger.KeyWorkerID, w.id, logger.KeyJobID, job.ID}
			if err != nil {
				w.workerLogger.With(attrs...).Error("Job failed", "error", err)
			} else {
				w.workerLogger.With(attrs...).Debug("Job completed")
			}
		case <-w.quit:
			return
		}
	}
}
