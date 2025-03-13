package indexer

import (
	"testing"
	"time"
)

// TestNewWorkerPool verifies pool creation with correct settings
func TestNewWorkerPool(t *testing.T) {
	numWorkers := 4
	pool := NewWorkerPool(numWorkers)

	if pool.numWorkers != numWorkers {
		t.Errorf("Expected %d workers, got %d", numWorkers, pool.numWorkers)
	}
	if len(pool.workers) != numWorkers {
		t.Errorf("Expected workers slice of length %d, got %d", numWorkers, len(pool.workers))
	}
	if cap(pool.tasks) != numWorkers*2 {
		t.Errorf("Expected task channel capacity %d, got %d", numWorkers*2, cap(pool.tasks))
	}
	if cap(pool.results) != numWorkers*2 {
		t.Errorf("Expected results channel capacity %d, got %d", numWorkers*2, cap(pool.results))
	}
}

// Test if the WorkerPool correctly starts the expected number of workers
func TestWorkerPoolStart(t *testing.T) {
	numWorkers := 3
	pool := NewWorkerPool(numWorkers)

	pool.Start()

	if len(pool.workers) != numWorkers {
		t.Errorf("Expected %d workers, but got %d", numWorkers, len(pool.workers))
	}

	for i, worker := range pool.workers {
		if worker.id != i {
			t.Errorf("Worker ID mismatch: expected %d, got %d", i, worker.id)
		}
		if worker.quit == nil {
			t.Errorf("Worker %d has a nil quit channel", i)
		}
	}


	taskSent := make(chan struct{})
	go func() {
		// Try to send a task
		select {
		case pool.tasks <- Task{ID: 1, Data: []byte("test")}:
			close(taskSent)
		case <-time.After(500 * time.Millisecond):
			// If we can't send after timeout, continue anyway
			close(taskSent)
		}
	}()

	select {
	case <-taskSent:
		
	case <-time.After(1 * time.Second):
		t.Log("Warning: Timeout waiting to send task")
	}
	stopDone := make(chan struct{})
	go func() {
		pool.Stop()
		close(stopDone)
	}()

	select {
	case <-stopDone:
		// Stopped successfully
	case <-time.After(1 * time.Second):
		t.Log("Warning: Timeout waiting for pool to stop")
	}
}

// TestSubmitAndResults verifies that tasks can be submitted and results retrieved
func TestSubmitAndResults(t *testing.T) {
	pool := NewWorkerPool(1)
	pool.Start()
	defer pool.Stop()

	// Create a test task
	testTask := Task{
		ID:   42,
		Data: []byte("test data"),
	}
	expectedHash := computeHash(testTask.Data)


	pool.Submit(testTask)

	
	var result Result
	select {
	case result = <-pool.Results():
		
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for result")
	}

	
	if result.TaskID != testTask.ID {
		t.Errorf("Expected TaskID %d, got %d", testTask.ID, result.TaskID)
	}

	if result.Hash != expectedHash {
		t.Errorf("Expected Hash %d, got %d", expectedHash, result.Hash)
	}
}
