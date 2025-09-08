package main

import (
	"context"
	"log/slog"
	"math/rand"
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
			Level:    slog.LevelDebug,
			ColorMap: logger.DefaultColorMap},
	)

	slog.SetDefault(slog.New(logHandler))
	slog.Info("Logger initialized")

	pool := worker.NewPool(3, 0)
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

	for i := 0; i < 10; i++ {
		someValTheWorkerDoesNotKnow := rand.Intn(10) + 1
		ctx := context.Background()

		newJob := worker.NewJob(ctx, func() (any, error) {
			res := strutil.LoremSentences(someValTheWorkerDoesNotKnow)
			return res, nil
		}).WithRetry(3, 1000).
			WithTimeout(1 * time.Second)

		pool.Submit(newJob)
	}

	pool.Shutdown()

	resultsWg.Wait()

	slog.With(slog.Int("processed", processed), slog.Int("failed", failed)).Info("Finished")
}
