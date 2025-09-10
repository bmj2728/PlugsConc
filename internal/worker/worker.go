package worker

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"
)

// ctxKeyWorkerID is the context key used to store and retrieve the worker ID from a context.
const (
	ctxKeyWorkerID = ctxKey("worker_id")
)

// KeyWorkerID is a constant key used to associate a worker's unique ID with context or logging operations.
const (
	KeyWorkerID = "worker_id"
)

// WithWorkerID returns a new context based on the parent context with the worker ID stored as a value.
func WithWorkerID(parent context.Context, id int) context.Context {
	return context.WithValue(parent, ctxKeyWorkerID, id)
}

//// WorkerIDFromContext extracts the worker ID as an integer from the given context.
//func WorkerIDFromContext(ctx context.Context) int {
//	return ctx.Value(ctxKeyWorkerID).(int)
//}

// Worker represents a worker that processes jobs from the jobs channel and sends results
// to the results channel.
type Worker struct {
	id      int
	jobs    <-chan *Job
	results chan<- *JobResult
	metrics chan<- *MetricResult
	quit    chan struct{}
}

// NewWorker creates and initializes a new Worker with a unique ID, a channel of jobs to process,
// and a results channel.
func NewWorker(id int, jobs <-chan *Job,
	results chan<- *JobResult,
	quit chan struct{},
	metrics chan<- *MetricResult) *Worker {
	return &Worker{
		id:      id,
		jobs:    jobs,
		results: results,
		quit:    quit,
		metrics: metrics,
	}
}

// Start begins the worker's execution loop, processing jobs from the channel and sending results
// to the results channel.
func (w *Worker) Start() {
	slogWorkerID := slog.Int(KeyWorkerID, w.id)
	slog.With(slogWorkerID).Debug("Worker started")
	defer slog.With(slogWorkerID).Debug("Worker stopped")

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
				for attempt := 0; ; attempt++ {
					job.RetryCount = attempt

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
					if e == nil || attempt >= job.MaxRetries {
						job.SetFinishedAt()
						return v, e
					}

					// log retry
					slog.With(
						slogWorkerID,
						slog.String(KeyJobID, job.ID),
						slog.Int(KeyRetryCount, attempt+1),
					).Warn("Retrying job")

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
				slog.With(slogWorkerID, job.LogValue()).Warn("Worker terminated before sending result")
				return
			}

			attrs := []any{slogWorkerID, slog.String(KeyJobID, job.ID)}
			if err != nil {
				slog.With(attrs...).Error("Job failed", slog.Any("error", err))
			} else {
				slog.With(attrs...).Debug("Job completed")
			}
		case <-w.quit:
			return
		}
	}
}
