package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

const (
	ctxKeyWorkerID = ctxKey("worker_id")
)

// JobWorkerIDKey is a constant key used to associate a worker's unique ID with context or logging operations.
const (
	JobWorkerIDKey = "worker_id"
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
	quit    chan struct{}
}

// NewWorker creates and initializes a new Worker with a unique ID, a channel of jobs to process,
// and a results channel.
func NewWorker(id int, jobs <-chan *Job, results chan<- *JobResult, quit chan struct{}) *Worker {
	return &Worker{
		id:      id,
		jobs:    jobs,
		results: results,
		quit:    quit,
	}
}

// Start begins the worker's execution loop, processing jobs from the channel and sending results
// to the results channel.
func (w *Worker) Start() {
	slogWorkerID := slog.Int(JobWorkerIDKey, w.id)
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
						err = fmt.Errorf("panic: %v", r)
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
						return nil, job.Ctx.Err()
					default:
					}

					// execute the job
					v, e := job.Execute(job.Ctx)
					// if the job succeeded, or we've reached the max retries, return the result/error
					//  otherwise, retry the job with a delay between retries'
					if e == nil || attempt >= job.MaxRetries {
						return v, e
					}

					// log retry
					slog.With(
						slogWorkerID,
						slog.String(JobIDKey, job.ID),
						slog.Int(RetryCountKey, attempt+1),
					).Warn("Retrying job")

					// wait for the retry delay before continuing the loop
					if delay > 0 {
						t := time.NewTimer(delay)
						// if the job context is canceled, stop the timer and return immediately,
						//  otherwise, wait for the timer to expire
						select {
						case <-job.Ctx.Done():
							t.Stop()
							return nil, job.Ctx.Err()
						case <-t.C:
						}
					}
				}
			}()

			// Safely send the result or quit if the pool is terminated.
			select {
			case w.results <- NewJobResult(job, resultVal, err):
				// Result sent successfully.
			case <-w.quit:
				// Pool was terminated while trying to send the result.
				// Log that the result is being discarded and exit the worker.
				slog.With(slogWorkerID, job.LogValue()).Warn("Worker terminated before sending result")
				return
			}

			attrs := []any{slogWorkerID, slog.String(JobIDKey, job.ID)}
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
