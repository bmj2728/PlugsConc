package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Worker represents a worker that processes jobs from the jobs channel and sends results
// to the results channel.
type Worker struct {
	id      int
	jobs    <-chan *Job
	results chan<- *JobResult
}

// JobWorkerIDKey is a constant key used to associate a worker's unique ID with context or logging operations.
const (
	JobWorkerIDKey = "worker_id"
)

// NewWorker creates and initializes a new Worker with a unique ID, a channel of jobs to process,
// and a results channel.
func NewWorker(id int, jobs <-chan *Job, results chan<- *JobResult) *Worker {
	return &Worker{
		id:      id,
		jobs:    jobs,
		results: results,
	}
}

// Start begins the worker's execution loop, processing jobs from the channel and sending results
// to the results channel.
func (w *Worker) Start() {
	slogWorkerID := slog.Int(JobWorkerIDKey, w.id)
	slog.With(slogWorkerID).Debug("Worker started")
	defer slog.With(slogWorkerID).Debug("Worker stopped")

	for job := range w.jobs {
		// annotate job context (assuming ctxKey is your private typed key)
		job.Ctx = context.WithValue(job.Ctx, ctxKeyWorkerID, w.id)

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

				select {
				case <-job.Ctx.Done():
					return nil, job.Ctx.Err()
				default:
				}

				v, e := job.Execute(job.Ctx)
				if e == nil || attempt >= job.MaxRetries {
					val, err = v, e
					return
				}

				// log retry
				slog.With(
					slogWorkerID,
					slog.String(JobIDKey, job.ID),
					slog.Int(RetryCountKey, attempt+1),
				).Warn("Retrying job")

				if delay > 0 {
					t := time.NewTimer(delay)
					select {
					case <-job.Ctx.Done():
						t.Stop()
						return nil, job.Ctx.Err()
					case <-t.C:
					}
				}
			}
		}()

		w.results <- NewJobResult(job, resultVal, err)

		attrs := []any{slogWorkerID, slog.String(JobIDKey, job.ID)}
		if err != nil {
			slog.With(attrs...).Error("Job failed", slog.Any("error", err))
		} else {
			slog.With(attrs...).Debug("Job completed")
		}
	}
}
