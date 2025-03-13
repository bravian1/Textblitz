package indexer

import (
	"sync"
	"time"
)

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
		go worker.run(&p.wg)
	}
	
}


func (p *WorkerPool) Submit(task Task) {
	p.tasks <- task
}


func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

		
func (p *WorkerPool) Stop() {
	
	close(p.tasks)

	
	for _, w := range p.workers {
		if w != nil { 
			select {
			case w.quit <- true: 
			default: 
			}
		}
	}

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):

	}


	close(p.results)
}

// run is the main worker processing loop
func (w *Worker) run(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-w.tasks:
			if !ok {
				return
			}
			result := Result{
				TaskID: task.ID,
				Hash:   computeHash(task.Data),
			}
			w.results <- result

		case <-w.quit:
			return 
		}
	}
}

func computeHash(data []byte) int64 {
	var hash int64
	for _, b := range data {
		hash = hash*31 + int64(b)
	}
	return hash
}
