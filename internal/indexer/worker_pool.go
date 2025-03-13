package indexer

import "sync"

// Task represents a chunk of data to be processed by a worker
type Task struct {
	ID int
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
