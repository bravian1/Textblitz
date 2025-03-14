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
// TestMultipleTasksProcessing tests that multiple tasks are correctly processed
func TestMultipleTasksProcessing(t *testing.T) {
	// Create a worker pool with 4 workers
	pool := NewSimHashWorkerPool(4)
	pool.Start()
	defer pool.Stop()

	// Create multiple test tasks with different content
	testStrings := []string{
		"This is the first test string",
		"This is the second test string with different content",
		"A completely different third string for testing",
		"Fourth and final test string with unique content",
	}

	// Send all tasks to the pool
	for i, str := range testStrings {
		task := Task{
			ID:         i,
			Data:       []byte(str),
			Offset:     i * 100, // Simulating different offsets
			SourceFile: "test.txt",
		}
		pool.Submit(task)
	}

	// Collect all results
	results := make(map[int]SimHashResult)
	for i := 0; i < len(testStrings); i++ {
		select {
		case result := <-pool.Results():
			results[result.TaskID] = result
		case <-time.After(2 * time.Second):
			t.Fatalf("Timed out waiting for result %d", i)
		}
	}

	// Verify all tasks were processed
	if len(results) != len(testStrings) {
		t.Errorf("Expected %d results, got %d", len(testStrings), len(results))
	}

	// Verify each task result
	for i, str := range testStrings {
		result, ok := results[i]
		if !ok {
			t.Errorf("No result found for task ID %d", i)
			continue
		}

		if !bytes.Equal(result.Data, []byte(str)) {
			t.Errorf("Task %d: Data mismatch", i)
		}

		if result.Offset != i*100 {
			t.Errorf("Task %d: Expected offset %d, got %d", i, i*100, result.Offset)
		}
	}

	// Verify each task produced a different hash (since content is different)
	hashes := make(map[uint64]bool)
	for _, result := range results {
		if hashes[result.Hash] {
			t.Error("Duplicate hash found, expected all hashes to be unique")
		}
		hashes[result.Hash] = true
	}
}

// TestSimilarContent tests that similar content produces similar hashes
func TestSimilarContent(t *testing.T) {
	pool := NewSimHashWorkerPool(1)
	pool.Start()
	defer pool.Stop()

	// Create two similar strings (only one word difference)
	str1 := "The quick brown fox jumps over the lazy dog"
	str2 := "The quick brown fox jumps over the lazy cat"

	// Submit both tasks
	pool.Submit(Task{ID: 1, Data: []byte(str1), Offset: 0, SourceFile: "test.txt"})
	pool.Submit(Task{ID: 2, Data: []byte(str2), Offset: 100, SourceFile: "test.txt"})

	// Get results
	var hash1, hash2 uint64
	for i := 0; i < 2; i++ {
		select {
		case result := <-pool.Results():
			if result.TaskID == 1 {
				hash1 = result.Hash
			} else {
				hash2 = result.Hash
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Timed out waiting for result")
		}
	}

	// The hashes should be different but similar (small Hamming distance)
	if hash1 == hash2 {
		t.Error("Expected different hashes for similar but different content")
	}

	// Count the number of differing bits (Hamming distance)
	xor := hash1 ^ hash2
	hammingDistance := 0
	for xor != 0 {
		hammingDistance++
		xor &= xor - 1 // Clear the least significant bit set
	}

	// For slightly different content, we expect the Hamming distance to be small
	// compared to the total bit length (64 bits)
	if hammingDistance > 20 {
		t.Errorf("Expected small Hamming distance for similar content, got %d", hammingDistance)
	}
}