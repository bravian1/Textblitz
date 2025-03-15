package internals

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

// Lookup handles the entire lookup workflow: load index, search hash, print results
// 1.Load index from file
// 2.Perform lookup
// 3.Print results
func LookUp(input_file string, simHash string, threshold int) error {
	indexmap, err := Load(input_file)
	if err != nil {
		return fmt.Errorf("Error loading index: %v\n", err)
	}

	fmt.Printf("Parsing simHash: %s\n", simHash)

	queryHash, err := strconv.ParseUint(simHash, 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid simHash format: %v", err)
	}

	fmt.Printf("Parsed queryHash: %d\n", queryHash)

	var matchedEntries []IndexEntry
	for key, entries := range indexmap {
		candidateHash, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			continue
		}
		if hammingDistance(queryHash, candidateHash) <= threshold {
			matchedEntries = append(matchedEntries, entries...)
		}
	}

	if len(matchedEntries) == 0 {
		return fmt.Errorf("No fuzzy matches found for SimHash: %s with threshold %d\n", simHash, threshold)
	}

	LookUpOutput(simHash, matchedEntries)
	return nil
}

// Calculates the number of differing bits between two 64-bit hashes.
func hammingDistance(a, b uint64) int {
	diff := a ^ b
	count := 0
	for diff != 0 {
		count++
		diff &= diff - 1
	}
	return count
}
