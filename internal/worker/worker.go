package worker

import (
	"context"
	"log/slog"
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
	for job := range w.jobs {
		job.Ctx = context.WithValue(job.Ctx, JobWorkerIDKey, w.id)

		// actual job execution
		result, err := job.Execute()

		if job.Cancel != nil {
			job.Cancel()
		} else if job.CancelWithCause != nil {
			job.CancelWithCause(job.Ctx.Err())
		}

		// send result to results channel
		w.results <- NewJobResult(job, result, err)

		slogJobID := slog.String(JobIDKey, job.ID)
		slog.With(slogWorkerID, slogJobID).Debug("Job completed")
	}
}
