package internals

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	idx "github.com/bravian1/Textblitz/internals/indexer"
)

// IndexFile processes a file, chunks it, computes simhashes for each chunk,
// and saves the indexed data to a file. It uses a worker pool for parallel processing.
func IndexFile(filename string, chunkSize int, numWorkers int, outputFile string) error {
	// Create an index manager to store our results
	indexManager := NewIndexManager()

	// Use the Chunk function to read and chunk the file
	chunks, err := idx.Chunk(filename, chunkSize)
	if err != nil {
		return fmt.Errorf("failed to chunk file: %w", err)
	}

	// Create a worker pool for parallel processing
	pool := idx.NewSimHashWorkerPool(numWorkers)
	pool.Start()

	// Create a channel to collect results that's large enough to prevent blocking
	resultChan := make(chan bool)

	// Process results in a background goroutine
	go func() {
		// Process all results from the worker pool
		for result := range pool.Results() {
			// Create an index entry for this chunk
			entry := IndexEntry{
				OriginalFile:    result.SourceFile,
				Size:            len(result.Data),
				Position:        result.Offset,
				AssociatedWords: extractKeywords(string(result.Data), 5),
			}

			// Add the entry to our index, keyed by its simhash
			if err := indexManager.Add(strconv.FormatUint(result.Hash, 10), entry); err != nil {
				fmt.Printf("Warning: failed to add entry to index: %v\n", err)
			}
		}

		// Signal that we're done processing results
		resultChan <- true
	}()

	// Submit all chunks to the worker pool
	for i, chunk := range chunks {
		pool.Submit(idx.Task{
			ID:         i,
			Data:       chunk,
			Offset:     i * chunkSize,
			SourceFile: filename,
		})
	}

	// Stop the worker pool (this will close the tasks channel)
	pool.Stop()

	// Wait for result processing to complete
	<-resultChan

	// If output file wasn't specified, generate one based on the input filename
	if outputFile == "" {
		outputFile = outputFilename(filename)
	}

	fmt.Printf("Saving index to: %s\n", outputFile)

	// Save the index to disk
	if err := indexManager.Save(outputFile); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	return nil
}

// outputFilename generates the index filename from the input filename
func outputFilename(inputFile string) string {
	extension := filepath.Ext(inputFile)
	baseName := strings.TrimSuffix(inputFile, extension)
	return baseName + ".idx"
}

// extractKeywords extracts a specified number of words from text for context
func extractKeywords(text string, count int) []string {
	words := strings.Fields(text)
	if len(words) <= count {
		return words
	}
	return words[:count]
}
