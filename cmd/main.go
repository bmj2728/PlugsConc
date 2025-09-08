package main

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"os"
	"sync"
	"time"

	"PlugsConc/internal/logger"
	"PlugsConc/internal/worker"

	"github.com/bmj2728/utils/pkg/strutil"
)

func main() {

	logHandler := logger.New(os.Stdout,
		&logger.Options{
			Level:    slog.LevelInfo,
			ColorMap: logger.DefaultColorMap},
	)

	slog.SetDefault(slog.New(logHandler))
	slog.Info("Logger initialized")

	pool := worker.NewPool(24, 0)
	processed := 0
	failed := 0

	go pool.Run()

	var resultsWg sync.WaitGroup
	resultsWg.Add(1)

	go func() {
		defer resultsWg.Done()
		for result := range pool.Results() {
			if result.Err == nil {
				processed++
			} else {
				failed++
			}
		}
	}()

	var jobsWithRetry []*worker.Job

	for i := 0; i < 250000; i++ {
		someValTheWorkerDoesNotKnow := rand.IntN(5) + 1
		ctx := context.Background()

		newJob := worker.NewJob(ctx, func(ctx context.Context) (any, error) {
			res := strutil.NewLoremSentences(someValTheWorkerDoesNotKnow)
			return res, nil
		}).WithRetry(3, 1000).
			WithTimeout(1 * time.Second)

		jobsWithRetry = append(jobsWithRetry, newJob)
	}

	s, f, err := pool.SubmitBatch(jobsWithRetry)
	if err != nil {
		slog.With(slog.Int("success", s), slog.Int("failed", f), slog.Any("err", err)).Info("Finished With Err")
	}

	pool.Shutdown()

	resultsWg.Wait()

	slog.With(slog.Int("success", processed),
		slog.Int("failed", failed),
		slog.Any("pool", pool.LogValue())).Info("Finished No Err")
}
