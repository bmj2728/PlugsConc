package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/bmj2728/PlugsConc/internal/logger"

	"github.com/bmj2728/utils/pkg/strutil"
)

// WorkUnit defines a function type that performs a unit of work and returns a value of any type.
type WorkUnit func(ctx context.Context) (any, error)

// Job represents a unit of work with an associated unique identifier and an executable function.
type Job struct {
	ID              string
	Metrics         *JobMetrics
	Execute         WorkUnit
	Ctx             context.Context
	Cancel          context.CancelFunc
	CancelWithCause context.CancelCauseFunc // only available if the job was created with WithCancelCause
	MaxRetries      int
	RetryDelay      int
}

// NewJob creates and initializes a new Job instance with a unique ID and the provided execution logic.
func NewJob(ctx context.Context, execute WorkUnit) *Job {
	uuid := strutil.GenerateUUIDV7()
	updatedCtx := WithJobID(ctx, uuid)
	return &Job{
		ID:      uuid,
		Execute: execute,
		Ctx:     updatedCtx,
		Metrics: NewJobMetrics(),
	}
}

// WithRetry configures the job with a maximum number of retries and a delay between retries in milliseconds.
func (j *Job) WithRetry(maxRetries int, retryDelay int) *Job {
	j.MaxRetries = maxRetries
	j.RetryDelay = retryDelay
	j.Ctx = context.WithValue(j.Ctx, ctxKeyMaxRetries, maxRetries)
	j.Ctx = context.WithValue(j.Ctx, ctxKeyRetryDelay, retryDelay)
	j.Ctx = context.WithValue(j.Ctx, ctxKeyRetryCount, 0)
	return j
}

// WithCancel creates a derived context with a cancel function for the current job and updates the job's context.
func (j *Job) WithCancel() *Job {
	updated, cancel := context.WithCancel(j.Ctx)
	j.Ctx = updated
	j.Cancel = cancel
	return j
}

// WithCancelCause updates the Job's context to include a cancelable context with a cause and returns the
// updated Job and a CancelCauseFunc.
func (j *Job) WithCancelCause() *Job {
	updated, cancel := context.WithCancelCause(j.Ctx)
	j.Ctx = updated
	j.CancelWithCause = cancel
	return j
}

// WithTimeout sets a timeout duration for the job's context and returns the updated job and a CancelFunc to cancel it.
func (j *Job) WithTimeout(timeout time.Duration) *Job {
	updated, cancel := context.WithTimeout(j.Ctx, timeout)
	j.Ctx = updated
	j.Cancel = cancel
	return j
}

// WithTimeoutCause sets a timeout and an associated cancellation cause for the Job's context, returning the
// updated Job and the cancel function.
func (j *Job) WithTimeoutCause(timeout time.Duration, cause error) *Job {
	updated, cancel := context.WithTimeoutCause(j.Ctx, timeout, cause)
	j.Ctx = updated
	j.Cancel = cancel
	return j
}

// WithDeadline sets a deadline on the Job's context and returns the Job and a CancelFunc to release resources.
func (j *Job) WithDeadline(deadline time.Time) *Job {
	updated, cancel := context.WithDeadline(j.Ctx, deadline)
	j.Ctx = updated
	j.Cancel = cancel
	return j
}

// WithDeadlineCause sets a deadline with a custom cancellation cause for the job's context and returns the
// updated job and cancel function.
func (j *Job) WithDeadlineCause(deadline time.Time, cause error) *Job {
	updated, cancel := context.WithDeadlineCause(j.Ctx, deadline, cause)
	j.Ctx = updated
	j.Cancel = cancel
	return j
}

// SetSubmittedAt updates the job's SubmittedAt field with the current time and stores it in the job's context.
func (j *Job) SetSubmittedAt() {
	j.Metrics.SubmittedAt = time.Now()
	j.Ctx = context.WithValue(j.Ctx, ctxKeyJobSubmittedAt, j.Metrics.SubmittedAt)
}

// SetStartedAt updates the Job's StartedAt timestamp and adds it to the Job's context as ctxKeyJobStartedAt.
func (j *Job) SetStartedAt() {
	j.Metrics.StartedAt = time.Now()
	j.Ctx = context.WithValue(j.Ctx, ctxKeyJobStartedAt, time.Now())
}

// SetFinishedAt sets the job's `FinishedAt` time to the current time, calculates the duration, and updates the context.
func (j *Job) SetFinishedAt() {
	j.Metrics.FinishedAt = time.Now()
	j.Ctx = context.WithValue(j.Ctx, ctxKeyJobFinishedAt, time.Now())
	j.Metrics.Duration = j.Metrics.FinishedAt.Sub(j.Metrics.StartedAt)
	j.Ctx = context.WithValue(j.Ctx, ctxKeyJobDuration, j.Metrics.Duration)
}

// JobResult represents the outcome of an operation with its associated JobID, result value, and any error encountered.
type JobResult struct {
	JobID    string
	WorkerID int
	Ctx      context.Context
	Metrics  *JobMetrics
	Value    any
	Err      error
}

// NewJobResult creates a new JobResult instance, copying the job's metrics and associating it with a specific worker.
func NewJobResult(job *Job, workerID int, value any, err error) *JobResult {
	return &JobResult{
		JobID:    job.ID,
		WorkerID: workerID,
		Ctx:      job.Ctx,
		Metrics:  job.Metrics,
		Value:    value,
		Err:      err,
	}
}

// LogValue returns a structured slog.Value representing key metadata and timing information for the Job.
func (j *Job) LogValue() slog.Value {
	return slog.GroupValue(slog.String(logger.KeyJobID, j.ID),
		slog.Int(logger.KeyMaxRetries, j.MaxRetries),
		slog.Int(logger.KeyRetryDelay, j.RetryDelay),
		slog.Any(logger.KeyJobMetrics, j.Metrics.LogValue()))
}

// LogValue returns a structured slog.Value containing detailed information about the job,
// such as timestamps and results.
func (jr *JobResult) LogValue() slog.Value {
	return slog.GroupValue(slog.String(logger.KeyJobID, jr.JobID),
		slog.Int(logger.KeyWorkerID, jr.WorkerID),
		slog.Any(logger.KeyJobValue, jr.Value),
		slog.Any(logger.KeyJobError, jr.Err),
		slog.Any(logger.KeyJobMetrics, jr.Metrics.LogValue()))
}
