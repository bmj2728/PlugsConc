package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/bmj2728/utils/pkg/strutil"
)

type ctxKey string

const (
	ctxKeyJobID          = ctxKey(KeyJobID)
	ctxKeyMaxRetries     = ctxKey(KeyMaxRetries)
	ctxKeyRetryDelay     = ctxKey(KeyRetryDelay)
	ctxKeyRetryCount     = ctxKey(KeyRetryCount)
	ctxKeyJobSubmittedAt = ctxKey(KeyJobSubmittedAt)
	ctxKeyJobStartedAt   = ctxKey(KeyJobStartedAt)
	ctxKeyJobFinishedAt  = ctxKey(KeyJobFinishedAt)
	ctxKeyJobDuration    = ctxKey(KeyJobDuration)
)

// WithJobID returns a copy of the parent context with the specified job ID added as a value.
func WithJobID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, ctxKeyJobID, id)
}

// JobIDFromContext retrieves the job ID from the given context. It assumes the context contains a value for
// the "job_id" key.
//func JobIDFromContext(ctx context.Context) string {
//	return ctx.Value(ctxKeyJobID).(string)
//}

const (
	KeyJobID          = "job_id"
	KeyMaxRetries     = "max_retries"
	KeyRetryDelay     = "retry_delay"
	KeyRetryCount     = "retry_count"
	KeyJobSubmittedAt = "submitted_at"
	KeyJobStartedAt   = "started_at"
	KeyJobFinishedAt  = "finished_at"
	KeyJobDuration    = "job_duration_seconds"
)

// WorkUnit defines a function type that performs a unit of work and returns a value of any type.
type WorkUnit func(ctx context.Context) (any, error)

// Job represents a unit of work with an associated unique identifier and an executable function.
type Job struct {
	ID              string
	SubmittedAt     time.Time
	StartedAt       time.Time
	FinishedAt      time.Time
	Duration        time.Duration
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
	updatedCtx := WithJobID(ctx, uuid)
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

func (j *Job) SetSubmittedAt() {
	j.SubmittedAt = time.Now()
	j.Ctx = context.WithValue(j.Ctx, ctxKeyJobSubmittedAt, j.SubmittedAt)
}

func (j *Job) SetStartedAt() {
	j.StartedAt = time.Now()
	j.Ctx = context.WithValue(j.Ctx, ctxKeyJobStartedAt, time.Now())
}

func (j *Job) SetFinishedAt() {
	j.FinishedAt = time.Now()
	j.Ctx = context.WithValue(j.Ctx, ctxKeyJobFinishedAt, time.Now())
	j.Duration = j.FinishedAt.Sub(j.StartedAt)
	j.Ctx = context.WithValue(j.Ctx, ctxKeyJobDuration, j.Duration)
}

// JobResult represents the outcome of an operation with its associated JobID, result value, and any error encountered.
type JobResult struct {
	JobID       string
	WorkerID    int
	Ctx         context.Context
	SubmittedAt time.Time
	StartedAt   time.Time
	FinishedAt  time.Time
	Duration    time.Duration
	Retries     int
	Value       any
	Err         error
}

// NewJobResult creates and returns a pointer to a JobResult containing the provided jobID, value, and error.
func NewJobResult(job *Job, workerID int, value any, err error) *JobResult {
	return &JobResult{
		JobID:       job.ID,
		WorkerID:    workerID,
		Ctx:         job.Ctx,
		SubmittedAt: job.SubmittedAt,
		StartedAt:   job.StartedAt,
		FinishedAt:  job.FinishedAt,
		Duration:    job.Duration,
		Retries:     job.RetryCount,
		Value:       value,
		Err:         err,
	}
}

func (j *Job) LogValue() slog.Value {
	return slog.GroupValue(slog.String(KeyJobID, j.ID),
		slog.Int(KeyMaxRetries, j.MaxRetries),
		slog.Int(KeyRetryDelay, j.RetryDelay),
		slog.Int(KeyRetryCount, j.RetryCount),
		slog.Time(KeyJobSubmittedAt, j.SubmittedAt),
		slog.Time(KeyJobStartedAt, j.StartedAt),
		slog.Time(KeyJobFinishedAt, j.FinishedAt),
		slog.Duration(KeyJobDuration, j.Duration),
	)
}

func (jr *JobResult) LogValue() slog.Value {
	return slog.GroupValue(slog.String(KeyJobID, jr.JobID),
		slog.Int(KeyWorkerID, jr.WorkerID),
		slog.Time(KeyJobSubmittedAt, jr.SubmittedAt),
		slog.Time(KeyJobStartedAt, jr.StartedAt),
		slog.Time(KeyJobFinishedAt, jr.FinishedAt),
		slog.Duration(KeyJobDuration, jr.Duration),
		slog.Int(KeyRetryCount, jr.Retries),
		slog.Any("value", jr.Value),
		slog.Any("error", jr.Err))
}
