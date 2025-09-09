package main

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"PlugsConc/internal/logger"
	"PlugsConc/internal/worker"

	"github.com/bmj2728/utils/pkg/strutil"
)

func main() {

	logHandler := logger.New(os.Stdout,
		&logger.Options{
			Level:     slog.LevelInfo,
			AddSource: true,
			ColorMap:  logger.DefaultColorMap},
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

	common := "With great power comes great responsibility"
	correct := "With great power there must also come great responsibility"

	for i := 0; i < 25; i++ {
		ctx := context.Background()
		newJob := worker.NewJob(ctx, func(ctx context.Context) (any, error) {
			res := strutil.CosineSimilarity(common, correct, 0)
			return res, nil
		}).WithRetry(3, 1000)

		jobsWithRetry = append(jobsWithRetry, newJob)
	}

	s, f, err := pool.SubmitBatch(jobsWithRetry)
	if err != nil {
		slog.With(slog.Int("submitted", s), slog.Int("failed", f), slog.Any("err", err)).Info("Batch Submission Finished")
	}

	pool.Shutdown()

	resultsWg.Wait()

	slog.With(slog.Int("success", processed),
		slog.Int("failed", failed),
		slog.Any("pool", pool.LogValue())).Info("Finished")

}
