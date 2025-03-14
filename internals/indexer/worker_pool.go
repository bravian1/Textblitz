package indexer

import (
	"sync"

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

func (p *WorkerPool) Stop() {
	close(p.tasks) // No more tasks

	// Wait for all workers to finish
	p.wg.Wait()

	// Close the results channel
	close(p.results)
}

func (w *SimHashWorker) run() {
	defer w.wg.Done()

	for {
		select {
		case task, ok := <-w.tasks:
			if !ok {
				return // Channel closed
			}

			text := string(task.Data)

			hash := w.simhasher.Hash(text)

			// Send result
			result := SimHashResult{
				TaskID:     task.ID,
				Hash:       hash,
				Data:       task.Data,
				Offset:     task.Offset,
				SourceFile: task.SourceFile,
			}
			w.results <- result

		case <-w.quit:
			return
		}
	}
}
