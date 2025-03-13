package internals

import (
	"fmt"

	"github.com/bravian1/Textblitz/internals/indexer"
)

func SystemIntegration(filename string, chunkSize int, numWorkers int) error {
	// Use the Chunk function to read and chunk the file
	_, err := indexer.Chunk(filename, chunkSize)
	if err != nil {
		return fmt.Errorf("failed to chunk file: %w", err)
	}

	return nil
}
