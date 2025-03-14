package internals

import (
	"encoding/gob"
	"encoding/json"
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

func (im *indexManager) Save(outputFile string) error {
	// Save in binary gob format for efficient loading
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(im.index); err != nil {
		return fmt.Errorf("failed to encode index: %w", err)
	}

	// Also save as JSON for human readability
	// Create a JSON file in the same directory as the index file
	jsonFilePath := outputFile + ".json"
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		fmt.Printf("Warning: Could not create JSON index file: %v\n", err)
		return nil // Don't fail the whole operation if JSON export fails
	}
	defer jsonFile.Close()

	// Use the encoding/json package
	jsonEncoder := json.NewEncoder(jsonFile)
	jsonEncoder.SetIndent("", "  ")
	if err := jsonEncoder.Encode(im.index); err != nil {
		fmt.Printf("Warning: Could not encode JSON index: %v\n", err)
	} else {
		fmt.Printf("Created human-readable index: %s\n", jsonFilePath)
	}

	return nil
}
