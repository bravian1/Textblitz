package indexer

import (
    "testing"
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
