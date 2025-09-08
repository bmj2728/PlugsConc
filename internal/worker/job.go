package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/bmj2728/utils/pkg/strutil"
)

const (
	JobIDKey      = "job_id"
	MaxRetriesKey = "max_retries"
	RetryDelayKey = "retry_delay"
	RetryCountKey = "retry_count"
)

// WorkUnit defines a function type that performs a unit of work and returns a value of any type.
type WorkUnit func() (any, error)

// Job represents a unit of work with an associated unique identifier and an executable function.
type Job struct {
	ID              string
	Execute         WorkUnit
	Ctx             context.Context
	Cancel          context.CancelFunc
	CancelWithCause context.CancelCauseFunc // only available if the job was created with WithCancelCause
	MaxRetries      int
	RetryDelay      int
	RetryCount      int
}

// NewJob creates and initializes a new Job instance with a unique ID and the provided execution logic.
func NewJob(ctx context.Context, execute WorkUnit) *Job {
	uuid := strutil.GenerateUUIDV7()
	updatedCtx := context.WithValue(ctx, JobIDKey, uuid)
	return &Job{
		ID:      uuid,
		Execute: execute,
		Ctx:     updatedCtx,
	}
}

// WithRetry configures the job with a maximum number of retries and a delay between retries in milliseconds.
func (j *Job) WithRetry(maxRetries int, retryDelay int) *Job {
	j.MaxRetries = maxRetries
	j.RetryDelay = retryDelay
	j.Ctx = context.WithValue(j.Ctx, MaxRetriesKey, maxRetries)
	j.Ctx = context.WithValue(j.Ctx, RetryDelayKey, retryDelay)
	j.Ctx = context.WithValue(j.Ctx, RetryCountKey, 0)
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

// JobResult represents the outcome of an operation with its associated JobID, result value, and any error encountered.
type JobResult struct {
	JobID string
	Ctx   context.Context
	Value any
	Err   error
}

// NewJobResult creates and returns a pointer to a JobResult containing the provided jobID, value, and error.
func NewJobResult(job *Job, value any, err error) *JobResult {
	return &JobResult{
		JobID: job.ID,
		Ctx:   job.Ctx,
		Value: value,
		Err:   err,
	}
}

func (j *Job) LogValue() slog.Value {
	return slog.GroupValue(slog.String(JobIDKey, j.ID),
		slog.Int(MaxRetriesKey, j.MaxRetries),
		slog.Int(RetryDelayKey, j.RetryDelay),
		slog.Int(RetryCountKey, j.RetryCount),
	)
}
