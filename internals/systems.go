package internals

import (
	"fmt"
	idx "github.com/bravian1/Textblitz/internal/indexer"
)

func SystemIntegration(filename string, chunkSize int, numWorkers int) error {
	// Use the Chunk function to read and chunk the file
	chunks, err := idx.Chunk(filename, chunkSize)
	if err != nil {
		return fmt.Errorf("failed to chunk file: %w", err)
	}

	pool := idx.NewWorkerPool(numWorkers)
}
