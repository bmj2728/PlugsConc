package worker

import "sync"

// Pool represents a worker pool that processes jobs concurrently with a fixed number of workers.
type Pool struct {
	maxWorkers int
	jobs       chan *Job
	results    chan *JobResult
	wg         *sync.WaitGroup
}

// NewPool creates and initializes a new Pool with the specified number of workers
// and an optional buffer size for channels.
func NewPool(maxWorkers int, buffer int) *Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	var jobs chan *Job
	var results chan *JobResult
	if buffer < 1 {
		// create unbuffered channels
		jobs = make(chan *Job)
		results = make(chan *JobResult)
	} else {
		// create buffered channels
		jobs = make(chan *Job, buffer)
		results = make(chan *JobResult, buffer)
	}
	return &Pool{
		maxWorkers: maxWorkers,
		jobs:       jobs,
		results:    results,
		wg:         &sync.WaitGroup{},
	}
}

// Run initializes and starts all workers in the pool, ensuring they are ready to process jobs concurrently.
func (p *Pool) Run() {
	for i := 1; i <= p.maxWorkers; i++ {
		nw := NewWorker(i, p.jobs, p.results)
		p.wg.Add(1)
		go func() {
			defer p.wg.Done() // Signal completion when the goroutine exits
			nw.Start()
		}()
	}
}

// Submit adds a job to the pool's job queue.
func (p *Pool) Submit(job *Job) {
	p.jobs <- job
}

// SubmitBatch submits a batch of jobs to the worker pool for processing.
func (p *Pool) SubmitBatch(jobs []*Job) {
	for _, job := range jobs {
		p.Submit(job)
	}
}

// Shutdown gracefully stops the pool by closing the job queue, waiting for workers to complete,
// and closing the results channel.
func (p *Pool) Shutdown() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}

// Results returns a read-only channel to retrieve processed job results from the worker pool.
func (p *Pool) Results() <-chan *JobResult {
	return p.results
}
