package indexer

import (
	"os"
	"sync"
)

func Chunk(filename string, chunkSize int) ([][]byte, error) {
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
