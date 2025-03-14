package indexer

import (
	"sync"
	"time"
	"github.com/bravian1/Textblitz/simhash"
)


type Task struct {
	ID         int    
	Data       []byte
	Offset     int    
	SourceFile string 
}

type SimHashResult struct {
	TaskID     int    
	Hash       uint64 
	Data       []byte
	Offset     int    
	SourceFile string 
}


type SimHashWorker struct {
	id        int                
	tasks     chan Task          
	results   chan SimHashResult  
	quit      chan bool           
	wg        *sync.WaitGroup     
	simhasher *simhash.SimHashGen
}

type WorkerPool struct {
	workers    []*SimHashWorker   
	numWorkers int               
	tasks      chan Task          
	results    chan SimHashResult 
	wg         sync.WaitGroup    
}

// NewSimHashWorkerPool creates and initializes a new worker pool with the specified number of workers.
// It sets up the necessary channels and worker instances but does not start them.
//
// Parameters:
//   - numWorkers: The number of worker goroutines to create
//
// Returns:
//   - *WorkerPool: A new worker pool instance ready to be started
func NewSimHashWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		workers:    make([]*SimHashWorker, numWorkers),
		numWorkers: numWorkers,
		tasks:      make(chan Task, numWorkers*2),
		results:    make(chan SimHashResult, numWorkers*2),
	}
}

// Start initializes and starts all workers in the pool.
// Each worker runs in its own goroutine and processes tasks until stopped.
// The method creates a shared feature set for all workers to ensure consistent hashing.
func (p *WorkerPool) Start() {
	featureSet := simhash.NewWordFeatureSet()

	for i := range p.numWorkers {
		p.wg.Add(1)
		worker := &SimHashWorker{
			id:        i,
			tasks:     p.tasks,
			results:   p.results,
			quit:      make(chan bool),
			wg:        &p.wg,
			simhasher: simhash.NewSimHashGenerator(featureSet),
		}
		p.workers[i] = worker
		go worker.run() // Start worker goroutine
	}
}
func (p *WorkerPool) Results() <-chan SimHashResult {
	return p.results
}

// Stop gracefully shuts down the worker pool.
// It closes the tasks channel, waits for all workers to finish,
// and then closes the results channel.
func (p *WorkerPool) Stop() {
	close(p.tasks) // No more tasks

	// Wait for all workers to finish
	p.wg.Wait()

	// Close the results channel
	close(p.results)
}

