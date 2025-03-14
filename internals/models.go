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

// IndexManager handles all operations related to the index
type IndexManager struct {
	index IndexMap
}

// NewIndexManager creates a new index manager
func NewIndexManager() *IndexManager {
	return &IndexManager{
		index: make(IndexMap),
	}
}

// Load reads an index from disk using gob encoding
func (im *IndexManager) Load(inputFile string) error {
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

// Lookup searches for entries with the given simhash value
func (im *IndexManager) Lookup(simhash string) ([]IndexEntry, error) {
	entries, ok := im.index[simhash]
	if !ok {
		return nil, fmt.Errorf("no entries found for SimHash: %s", simhash)
	}
	return entries, nil
}

// Add adds a new entry to the index
func (im *IndexManager) Add(simhash string, entry IndexEntry) error {
	im.index[simhash] = append(im.index[simhash], entry)
	return nil
}

// Save writes the index to disk in both binary (gob) and JSON formats
func (im *IndexManager) Save(outputFile string) error {
	// Save in binary gob format for efficient loading
	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(im.index); err != nil {
		return fmt.Errorf("failed to encode index: %w", err)
	}

	//Also save as JSON for human readability
	jsonFilePath := outputFile + ".json"
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		fmt.Printf("Warning: Could not create JSON index file: %v\n", err)
		return nil
	}
	defer jsonFile.Close()

	jsonEncoder := json.NewEncoder(jsonFile)
	jsonEncoder.SetIndent("", "  ")
	if err := jsonEncoder.Encode(im.index); err != nil {
		fmt.Printf("Warning: Could not encode JSON index: %v\n", err)
	} else {
		fmt.Printf("Created human-readable index: %s\n", jsonFilePath)
	}

	return nil
}

// LookUpOutput formats and prints the lookup results
func LookUpOutput(simHash string, entries []IndexEntry) {
	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return
	}

	fmt.Println("\nLookup Complete!")
	fmt.Println("------------------------------------")

	for _, entry := range entries {
		fmt.Printf("| SimHash       : %s\n", simHash)
		fmt.Printf("| Original File : %s\n", entry.OriginalFile)
		fmt.Printf("| Position      : Byte %d\n", entry.Position)
		fmt.Printf("| Associated Words : \"%s\"\n", entry.AssociatedWords)
		fmt.Println("------------------------------------------------")
	}

	fmt.Println()
}
