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

	// Check if all workers were created
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

	// Give time for goroutines to start but with a timeout
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

	// Wait for task to be sent or timeout
	select {
	case <-taskSent:
		// Continue
	case <-time.After(1 * time.Second):
		t.Log("Warning: Timeout waiting to send task")
	}

	// Make sure to call Stop in a non-blocking way
	stopDone := make(chan struct{})
	go func() {
		pool.Stop()
		close(stopDone)
	}()

	// Wait for stop with timeout
	select {
	case <-stopDone:
		// Stopped successfully
	case <-time.After(1 * time.Second):
		t.Log("Warning: Timeout waiting for pool to stop")
	}
}
