package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bmj2728/PlugsConc/internal/logger"
)

// ctxKey is a custom string-based type used as keys for storing and retrieving values in context.
type ctxKey string

const (
	ctxWarningPrefix = "Context does not contain the key"
	// ctxKeyJobID is the context key for storing or retrieving a unique job identifier.
	ctxKeyJobID = ctxKey(logger.KeyJobID)
	// ctxKeyMaxRetries is the context key for storing or retrieving the maximum allowed retries for a job.
	ctxKeyMaxRetries = ctxKey(logger.KeyMaxRetries)
	// ctxKeyRetryDelay is the context key for storing or retrieving the delay duration before retrying a job.
	ctxKeyRetryDelay = ctxKey(logger.KeyRetryDelay)
	// ctxKeyRetryCount is the context key for storing or retrieving the current retry count of a job.
	ctxKeyRetryCount = ctxKey(logger.KeyRetryCount)
	// ctxKeyJobSubmittedAt is the context key for storing or retrieving the job submission timestamp.
	ctxKeyJobSubmittedAt = ctxKey(logger.KeyJobSubmittedAt)
	// ctxKeyJobStartedAt is the context key for storing or retrieving the job start timestamp.
	ctxKeyJobStartedAt = ctxKey(logger.KeyJobStartedAt)
	// ctxKeyJobFinishedAt is the context key for storing or retrieving the job completion timestamp.
	ctxKeyJobFinishedAt = ctxKey(logger.KeyJobFinishedAt)
	// ctxKeyJobDuration is the context key for storing or retrieving the job's total execution duration in seconds.
	ctxKeyJobDuration = ctxKey(logger.KeyJobDuration)
	// ctxKeyWorkerCount is a context key for tracking the number of workers in a pool.
	ctxKeyWorkerCount = ctxKey(logger.KeyWorkerCount)
	// ctxKeySubmittedJobs is a context key for tracking the total number of submitted jobs.
	ctxKeySubmittedJobs = ctxKey(logger.KeySubmittedJobs)
	// ctxKeyFailedSubmissions is a context key for tracking the count of job submission failures.
	ctxKeyFailedSubmissions = ctxKey(logger.KeyFailedSubmissions)
	// ctxKeyPoolStartedAt is a context key for storing the pool's start time.
	ctxKeyPoolStartedAt = ctxKey(logger.KeyPoolStartedAt)
	// ctxKeyPoolStoppedAt is a context key for storing the pool's stop time.
	ctxKeyPoolStoppedAt = ctxKey(logger.KeyPoolStoppedAt)
	// ctxKeyPoolCompletedAt is a context key for storing the pool's completion time.
	ctxKeyPoolCompletedAt = ctxKey(logger.KeyPoolCompletedAt)
	// ctxKeyPoolDuration is a context key for tracking the total duration the pool was active.
	ctxKeyPoolDuration = ctxKey(logger.KeyPoolDuration)
	// ctxKeyPoolClosed is a context key for indicating whether the pool has been closed.
	ctxKeyPoolClosed = ctxKey(logger.KeyPoolClosed)
	// ctxKeySuccessfulJobs is a context key for tracking the number of successfully completed jobs.
	ctxKeySuccessfulJobs = ctxKey(logger.KeySuccessfulJobs)
	// ctxKeyFailedJobs is a context key for tracking the number of failed jobs.
	ctxKeyFailedJobs = ctxKey(logger.KeyFailedJobs)
	// ctxKeyWorkerID is the context key used to store and retrieve the worker ID from a context.
	ctxKeyWorkerID = ctxKey("worker_id")
)

// WithJobID returns a copy of the parent context with the specified job ID added as a value.
func WithJobID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, ctxKeyJobID, id)
}

// WithWorkerID returns a new context based on the parent context with the worker ID stored as a value.
func WithWorkerID(parent context.Context, id int) context.Context {
	return context.WithValue(parent, ctxKeyWorkerID, id)
}

// JobIDFromCtx retrieves the job ID from the given context. It assumes the context contains a value for
// the "job_id" key.
func JobIDFromCtx(ctx context.Context) string {
	val, ok := ctx.Value(ctxKeyJobID).(string)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyJobID))
		return ""
	}
	return val
}

// MaxRetriesFromCtx retrieves the maximum retry count from the provided context.
// Returns 0 if the value is not present or if an invalid value is encountered.
func MaxRetriesFromCtx(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeyMaxRetries).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyMaxRetries))
		return 0
	}
	return val
}

// RetryDelayFromCtx retrieves the retry delay duration stored in the given context as an integer.
// Returns 0 if the key is not present or its value is not an integer.
func RetryDelayFromCtx(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeyRetryDelay).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyRetryDelay))
		return 0
	}
	return val
}

// RetryCountFromCtx retrieves the retry count from the provided context, returning 0 if the key is
// not present or invalid.
func RetryCountFromCtx(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeyRetryCount).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyRetryCount))
		return 0
	}
	return val
}

// JobSubmittedAtFromCtx retrieves the job submission timestamp from the provided context using a predefined key.
// If the key is absent or the value cannot be cast to time.Time, it logs a warning and
// returns the zero value of time.Time.
func JobSubmittedAtFromCtx(ctx context.Context) time.Time {
	val, ok := ctx.Value(ctxKeyJobSubmittedAt).(time.Time)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyJobSubmittedAt))
		return time.Time{}
	}
	return val
}

// JobStartedAtFromCtx retrieves the job start timestamp stored in the context using the key ctxKeyJobStartedAt.
// If the key is not present or the value is not a time.Time, it logs a warning and returns the zero value of time.Time.
func JobStartedAtFromCtx(ctx context.Context) time.Time {
	val, ok := ctx.Value(ctxKeyJobStartedAt).(time.Time)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyJobStartedAt))
		return time.Time{}
	}
	return val
}

// JobFinishedAtFromCtx retrieves the job completion timestamp from the provided context using a specific context key.
// Returns an empty time.Time value if the key is not present or its type assertion fails.
// Logs a warning using slog if the key is missing or invalid.
func JobFinishedAtFromCtx(ctx context.Context) time.Time {
	val, ok := ctx.Value(ctxKeyJobFinishedAt).(time.Time)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyJobFinishedAt))
		return time.Time{}
	}
	return val
}

// JobDurationSecondsFromCtx retrieves the job's execution duration from the context or returns zero if
// unavailable or invalid.
func JobDurationSecondsFromCtx(ctx context.Context) time.Duration {
	val, ok := ctx.Value(ctxKeyJobDuration).(time.Duration)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyJobDuration))
		return 0
	}
	return val
}

// WorkerCountFromCtx retrieves the worker count from the provided context using a predefined context key.
// If the key is not present or its value is not an int, a warning is logged, and the function returns 0.
func WorkerCountFromCtx(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeyWorkerCount).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyWorkerCount))
		return 0
	}
	return val
}

// WorkerIDFromContext retrieves the worker ID from the provided context.
// Returns 0 and logs a warning if the key is not present or the value is not an integer.
func WorkerIDFromContext(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeyWorkerID).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyWorkerID))
		return 0
	}
	return val
}

// SubmittedJobsFromCtx retrieves the total number of submitted jobs from the provided context.
// Returns 0 if the key is missing or the value is not an integer.
func SubmittedJobsFromCtx(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeySubmittedJobs).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeySubmittedJobs))
		return 0
	}
	return val
}

// FailedSubmissionsFromCtx retrieves the count of failed submissions from the context or returns 0 if unavailable.
func FailedSubmissionsFromCtx(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeyFailedSubmissions).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyFailedSubmissions))
		return 0
	}
	return val
}

// PoolStartedAtFromCtx retrieves the pool's start time from the provided context using ctxKeyPoolStartedAt.
// If the key is not found, it logs a warning and returns the zero value of time.Time.
func PoolStartedAtFromCtx(ctx context.Context) time.Time {
	val, ok := ctx.Value(ctxKeyPoolStartedAt).(time.Time)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyPoolStartedAt))
		return time.Time{}
	}
	return val
}

// PoolStoppedAtFromCtx retrieves a time.Time value associated with ctxKeyPoolStoppedAt from the context.
// It logs a warning and returns the zero value of time.Time if the key is not found or
// if the value is not of time.Time type.
func PoolStoppedAtFromCtx(ctx context.Context) time.Time {
	val, ok := ctx.Value(ctxKeyPoolStoppedAt).(time.Time)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyPoolStoppedAt))
		return time.Time{}
	}
	return val
}

// PoolCompletedAtFromCtx retrieves the pool's completion time from the context using the ctxKeyPoolCompletedAt key.
// It returns a zero time if the key is missing or the value is not of type time.Time.
// Logs a warning if the context does not contain the expected key or value type mismatch occurs.
func PoolCompletedAtFromCtx(ctx context.Context) time.Time {
	val, ok := ctx.Value(ctxKeyPoolCompletedAt).(time.Time)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyPoolCompletedAt))
		return time.Time{}
	}
	return val
}

// PoolClosedFromCtx checks the context for the `ctxKeyPoolClosed` key and returns its value
// if it's a boolean, or false otherwise.
func PoolClosedFromCtx(ctx context.Context) bool {
	val, ok := ctx.Value(ctxKeyPoolClosed).(bool)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyPoolClosed))
		return false
	}
	return val
}

// PoolDurationFromCtx extracts the pool duration from the given context using the ctxKeyPoolDuration key.
// If the key is not present or the value is not of type time.Duration, it logs a warning and returns 0.
func PoolDurationFromCtx(ctx context.Context) time.Duration {
	val, ok := ctx.Value(ctxKeyPoolDuration).(time.Duration)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyPoolDuration))
		return 0
	}
	return val
}

// SuccessfulJobsFromCtx retrieves the number of successfully completed jobs from the provided context.
// If the context does not contain a valid value for the key, logs a warning and returns 0.
func SuccessfulJobsFromCtx(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeySuccessfulJobs).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeySuccessfulJobs))
		return 0
	}
	return val
}

// FailedJobsFromCtx retrieves the count of failed jobs from the provided context.
// It returns 0 if the key is missing or the value is not an integer.
func FailedJobsFromCtx(ctx context.Context) int {
	val, ok := ctx.Value(ctxKeyFailedJobs).(int)
	if !ok {
		slog.Warn(fmt.Sprintf("%s %q", ctxWarningPrefix, ctxKeyFailedJobs))
		return 0
	}
	return val
}
