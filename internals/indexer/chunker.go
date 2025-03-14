package indexer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"code.sajari.com/docconv"
)

// Chunk divides a file into chunks of specified size
// Supports .txt, .pdf, and .docx files
func Chunk(filename string, chunkSize int) ([][]byte, error) {
	// Get file extension
	ext := strings.ToLower(filepath.Ext(filename))
	// Process based on file type
	switch ext {
	case ".txt":
		return chunkFileWithGoroutines(filename, chunkSize)
	case ".pdf", ".docx", "xml" :
		data, err := extractTextFromDoc(filename)
		if err != nil {
			return nil, err
		}
		return chunkDataWithGoroutines(data, chunkSize)

	default:
		return nil, errors.New("unsupported file type: " + ext)
	}
}

// chunkFileWithGoroutines chunks a file using goroutines
func chunkFileWithGoroutines(filename string, chunkSize int) ([][]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	numChunks := int((fileSize + int64(chunkSize) - 1) / int64(chunkSize))

	chunks := make([][]byte, numChunks)
	var wg sync.WaitGroup
	wg.Add(numChunks)

	for i := 0; i < numChunks; i++ {
		go func(i int) {
			defer wg.Done()
			offset := int64(i * chunkSize)
			buf := make([]byte, chunkSize)
			n, _ := file.ReadAt(buf, offset)
			chunks[i] = buf[:n]
		}(i)
	}

	wg.Wait()
	return chunks, nil
}

// chunkDataWithGoroutines chunks a byte slice using goroutines
func chunkDataWithGoroutines(data []byte, chunkSize int) ([][]byte, error) {
	dataLen := len(data)
	numChunks := (dataLen + chunkSize - 1) / chunkSize

	chunks := make([][]byte, numChunks)
	var wg sync.WaitGroup
	wg.Add(numChunks)

	for i := 0; i < numChunks; i++ {
		go func(i int) {
			defer wg.Done()
			start := i * chunkSize
			end := start + chunkSize
			if end > dataLen {
				end = dataLen
			}
			chunks[i] = make([]byte, end-start)
			copy(chunks[i], data[start:end])
		}(i)
	}

	wg.Wait()
	return chunks, nil
}

// extractTextFromDoc extracts  text from pdf , docx or xml using sajari's docconv library
//
// takes a file path as input and returns the extracted text as slice of  bytes
func extractTextFromDoc(filename string) ([]byte, error) {
	result, err := docconv.ConvertPath(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document to text: %v", err)
	}
	return []byte(result.Body), nil
}
