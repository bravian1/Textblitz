package internals

import (
	"encoding/gob"
	"fmt"
	"os"
)

type IndexEntry struct {
	OriginalFile    string
	Size            int
	Position        int
	AssociatedWords []string
}

type IndexMap map[string][]IndexEntry

type IndexManager interface {
	Load(inputFile string) error

	Lookup(simhash string) ([]IndexEntry, error)

	Add(simhash string, entry IndexEntry) error

	Save(outputFile string) error
}

type indexManager struct {
	index IndexMap
}

// NewIndexManager creates a new index manager
func NewIndexManager() IndexManager {
	return &indexManager{
		index: make(IndexMap),
	}
}

// Load reads an index from disk using gob encoding
func (im *indexManager) Load(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open index file: %w", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&im.index); err != nil {
		return fmt.Errorf("failed to decode index: %w", err)
	}

	return nil
}



func (im *indexManager) Lookup(simhash string) ([]IndexEntry, error) {
	entries, found := im.index[simhash]
	if !found {
		return nil, fmt.Errorf("simhash %s not found", simhash)
	}
	return entries, nil
}


func (im *indexManager) Add(simhash string, entry IndexEntry) error {
	im.index[simhash] = append(im.index[simhash], entry)
	return nil
}
