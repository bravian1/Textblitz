package indexer

import (
	"bytes"
	"testing"
	"time"
)

// TestNewSimHashWorkerPool tests that a new worker pool is created with the correct configuration
func TestNewSimHashWorkerPool(t *testing.T) {
	numWorkers := 4
	pool := NewSimHashWorkerPool(numWorkers)

	if pool == nil {
		t.Fatal("NewSimHashWorkerPool returned nil")
	}

	if pool.numWorkers != numWorkers {
		t.Errorf("Expected numWorkers to be %d, got %d", numWorkers, pool.numWorkers)
	}

	if len(pool.workers) != numWorkers {
		t.Errorf("Expected workers slice to have length %d, got %d", numWorkers, len(pool.workers))
	}

	if cap(pool.tasks) != numWorkers*2 {
		t.Errorf("Expected tasks channel capacity to be %d, got %d", numWorkers*2, cap(pool.tasks))
	}

	if cap(pool.results) != numWorkers*2 {
		t.Errorf("Expected results channel capacity to be %d, got %d", numWorkers*2, cap(pool.results))
	}
}

// TestWorkerPoolBasicProcessing tests the basic flow of a worker pool processing tasks
func TestWorkerPoolBasicProcessing(t *testing.T) {
	// Create a worker pool with 2 workers
	pool := NewSimHashWorkerPool(2)
	pool.Start()
	defer pool.Stop()

	// Create a simple test task
	testData := []byte("This is a test string for simhash calculation")
	task := Task{
		ID:         1,
		Data:       testData,
		Offset:     0,
		SourceFile: "test.txt",
	}

	// Submit the task to the pool
	pool.Submit(task)

	// Get the result
	var result SimHashResult
	select {
	case result = <-pool.Results():
		// Got a result
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for result")
	}

	// Verify the result
	if result.TaskID != task.ID {
		t.Errorf("Expected TaskID %d, got %d", task.ID, result.TaskID)
	}

	if !bytes.Equal(result.Data, task.Data) {
		t.Error("Result data doesn't match task data")
	}

	if result.Offset != task.Offset {
		t.Errorf("Expected Offset %d, got %d", task.Offset, result.Offset)
	}

	if result.SourceFile != task.SourceFile {
		t.Errorf("Expected SourceFile %s, got %s", task.SourceFile, result.SourceFile)
	}

	if result.Hash == 0 {
		t.Error("Expected non-zero hash value")
	}
}
