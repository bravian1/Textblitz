package internals

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	idx "github.com/bravian1/Textblitz/internals/indexer"
)

func IndexFile(filename string, chunkSize int, numWorkers int) error {

	indexManager := NewIndexManager()

	chunks, err := idx.Chunk(filename, chunkSize)
	if err != nil {
		return fmt.Errorf("failed to chunk file: %w", err)
	}

	pool := idx.NewSimHashWorkerPool(numWorkers)
	pool.Start()

	for i, chunk := range chunks {
		pool.Submit(idx.Task{
			ID:         i,
			Data:       chunk,
			Offset:     i * chunkSize,
			SourceFile: filename,
		})
	}

	resultCount := 0

	results := make(chan idx.SimHashResult, len(chunks))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 0; i < len(chunks); i++ {
			result, ok := <-pool.Results()
			if !ok {
				break
			}
			results <- result
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {

		entry := IndexEntry{
			OriginalFile:    result.SourceFile,
			Size:            len(result.Data),
			Position:        result.Offset,
			AssociatedWords: extractKeywords(string(result.Data), 5),
		}

		if err := indexManager.Add(strconv.FormatUint(result.Hash, 10), entry); err != nil {
			pool.Stop()
			return fmt.Errorf("failed to add entry to index: %w", err)
		}

		resultCount++
	}

	pool.Stop()

	outputFile := outputFilename(filename)

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
