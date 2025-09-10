package worker

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"PlugsConc/internal/logger"
)

// ErrNoStart indicates that a required start time is missing.
// ErrNoComplete indicates that a required completion time is missing.
var (
	ErrNoStart    = errors.New("no start time exists")
	ErrNoComplete = errors.New("no complete time exists")
)

// PoolMetrics captures metrics about the lifecycle and performance of a thread pool during its runtime.
type PoolMetrics struct {
	mu                 sync.RWMutex  // mutex to allow threadsafe ops
	startedAt          time.Time     // when Run() was called
	stoppedAt          time.Time     // when Shutdown(), Stop(), or Terminate() were called
	completedAt        time.Time     // when last job was returned
	duration           time.Duration // from startedAt to completedAt
	submissions        int           // jobs submitted
	submissionFailures int           // jobs that were unable to be submitted
	succeeded          int           // jobs that completed successfully
	failed             int           // jobs that did not complete successfully
}

// NewPoolMetrics initializes a new instance of PoolMetrics with default values and a mutex for thread safety.
func NewPoolMetrics() *PoolMetrics {
	return &PoolMetrics{
		mu: sync.RWMutex{},
	}
}

// Started retrieves the timestamp when the pool was started. It is thread-safe.
func (pm *PoolMetrics) Started() time.Time {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.startedAt
}

// Stopped returns the time when the pool was stopped, ensuring thread-safe access to the stoppedAt field.
func (pm *PoolMetrics) Stopped() time.Time {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.stoppedAt
}

// Completed returns the time when the last job in the pool was completed in a thread-safe manner.
func (pm *PoolMetrics) Completed() time.Time {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.completedAt
}

// Duration returns the total duration between the pool's start and completion times in a thread-safe manner.
func (pm *PoolMetrics) Duration() time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.duration
}

// Submissions returns the total number of jobs submitted to the pool. It is a threadsafe operation.
func (pm *PoolMetrics) Submissions() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.submissions
}

// FailedSubmissions returns the total number of jobs that failed to be submitted to the pool.
func (pm *PoolMetrics) FailedSubmissions() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.submissionFailures
}

// SuccessfulJobs returns the number of jobs that have completed successfully in the pool. It locks for reading.
func (pm *PoolMetrics) SuccessfulJobs() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.succeeded
}

// FailedJobs returns the number of jobs that did not complete successfully.
func (pm *PoolMetrics) FailedJobs() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.failed
}

// SetStarted records the current time as the start time for the pool. It ensures thread safety using a mutex lock.
func (pm *PoolMetrics) SetStarted() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.startedAt = time.Now()
}

// SetStopped records the current time as the point when the pool was stopped, ensuring thread-safe access with a mutex.
func (pm *PoolMetrics) SetStopped() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.stoppedAt = time.Now()
}

// SetCompleted records the time when the last job in the pool was completed.
func (pm *PoolMetrics) SetCompleted() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.completedAt = time.Now()
}

// SetDuration calculates and sets the duration between the pool's start and completion times.
// Returns an error if either is unset.
func (pm *PoolMetrics) SetDuration() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if pm.startedAt.IsZero() {
		return ErrNoStart
	}
	if pm.completedAt.IsZero() {
		return ErrNoComplete
	}
	pm.duration = pm.completedAt.Sub(pm.startedAt)
	return nil
}

// RecordSubmission increments the count of successfully submitted jobs in a thread-safe manner.
func (pm *PoolMetrics) RecordSubmission() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.submissions++
}

// RecordFailedSubmission increments the count of failed job submissions within the pool metrics.
func (pm *PoolMetrics) RecordFailedSubmission() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.submissionFailures++
}

// RecordSuccessfulJob increments the count of jobs that completed successfully in a thread-safe manner.
func (pm *PoolMetrics) RecordSuccessfulJob() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.succeeded++
}

// RecordFailedJob increments the count of jobs that did not complete successfully in a thread-safe manner.
func (pm *PoolMetrics) RecordFailedJob() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.failed++
}

// LogValue returns a slog.Value representing the current state of PoolMetrics, including job counts and time attributes.
func (pm *PoolMetrics) LogValue() slog.Value {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	metrics := slog.GroupValue(slog.Int(logger.KeySubmittedJobs, pm.submissions),
		slog.Int(logger.KeyFailedSubmissions, pm.failed),
		slog.Int(logger.KeySuccessfulJobs, pm.succeeded),
		slog.Int(logger.KeyFailedJobs, pm.failed),
		slog.Time(logger.KeyPoolStartedAt, pm.startedAt),
		slog.Time(logger.KeyPoolStoppedAt, pm.stoppedAt),
		slog.Time(logger.KeyPoolCompletedAt, pm.completedAt),
		slog.Float64(logger.KeyPoolDuration, pm.duration.Seconds()))
	return metrics
}
