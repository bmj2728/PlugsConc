package worker

import (
	"errors"
	"log/slog"
	"sync"
	"time"
)

var (
	ErrNoStart    = errors.New("no start time exists")
	ErrNoComplete = errors.New("no complete time exists")
)

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

func NewPoolMetrics() *PoolMetrics {
	return &PoolMetrics{
		mu: sync.RWMutex{},
	}
}

func (pm *PoolMetrics) Started() time.Time {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.startedAt
}

func (pm *PoolMetrics) Stopped() time.Time {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.stoppedAt
}

func (pm *PoolMetrics) Completed() time.Time {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.completedAt
}

func (pm *PoolMetrics) Duration() time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.duration
}

func (pm *PoolMetrics) Submissions() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.submissions
}

func (pm *PoolMetrics) FailedSubmissions() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.submissionFailures
}

func (pm *PoolMetrics) SuccessfulJobs() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.succeeded
}

func (pm *PoolMetrics) FailedJobs() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.failed
}

func (pm *PoolMetrics) SetStarted() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.startedAt = time.Now()
}

func (pm *PoolMetrics) SetStopped() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.stoppedAt = time.Now()
}

func (pm *PoolMetrics) SetCompleted() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.completedAt = time.Now()
}

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

func (pm *PoolMetrics) RecordSubmission() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.submissions++
}

func (pm *PoolMetrics) RecordFailedSubmission() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.submissionFailures++
}

func (pm *PoolMetrics) RecordSuccessfulJob() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.succeeded++
}

func (pm *PoolMetrics) RecordFailedJob() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.failed++
}

func (pm *PoolMetrics) LogValue() slog.Value {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	metrics := slog.GroupValue(slog.Int(KeySubmittedJobs, pm.submissions),
		slog.Int(KeyFailedSubmissions, pm.failed),
		slog.Int(KeySuccessfulJobs, pm.succeeded),
		slog.Int(KeyFailedJobs, pm.failed),
		slog.Time(KeyPoolStartedAt, pm.startedAt),
		slog.Time(KeyPoolStoppedAt, pm.stoppedAt),
		slog.Time(KeyPoolCompletedAt, pm.completedAt),
		slog.Float64(KeyPoolDuration, pm.duration.Seconds()))
	return metrics
}
