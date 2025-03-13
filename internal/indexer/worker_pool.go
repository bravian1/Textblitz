package indexer

import "sync"

// Task represents a chunk of data to be processed by a worker
type Task struct {
	ID   int
	Data []byte
}

// Result contains the output of processing a task
type Result struct {
	TaskID int
	Hash   int64
}

// Worker represents a single worker goroutine
type Worker struct {
	id      int
	tasks   chan Task
	results chan Result
	quit    chan bool
}

// WorkerPool manages multiple workers for parallel processing
type WorkerPool struct {
	workers    []*Worker
	numWorkers int
	tasks      chan Task
	results    chan Result
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		workers:    make([]*Worker, numWorkers),
		numWorkers: numWorkers,
		tasks:      make(chan Task, numWorkers*2),   // Buffered channel
		results:    make(chan Result, numWorkers*2), // Buffered channel
	}
}
func (p *WorkerPool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		worker := &Worker{
			id:      i,
			tasks:   p.tasks,
			results: p.results,
			quit:    make(chan bool),
		}
		p.workers[i] = worker
		p.wg.Add(1)
		go worker.run(&p.wg) // Start in goroutine
	}
}
// Submit adds a new task to the pool
func (p *WorkerPool) Submit(task Task) {
    p.tasks <- task
}

// Results returns the channel for receiving results
func (p *WorkerPool) Results() <-chan Result {
    return p.results
}

// Stop gracefully shuts down the worker pool
func (p *WorkerPool) Stop() {
    close(p.tasks)          // No more tasks
    for _, w := range p.workers {
        w.quit <- true      // Signal workers
    }
    p.wg.Wait()            // Wait for completion
    close(p.results)       // Close results
}


// Start initializes and starts all workers in the pool
