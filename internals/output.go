package internals

import (
	"encoding/gob"
	"fmt"
	"os"
)

// save: writes the indexMap to  a file with gob
func Save(filename string, indexmap IndexMap) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(indexmap); err != nil {
		return err
	}
	return nil
}

// Formats the outputs for lookup. It takes a slice of IndexEntry as input and prints the formatted outputs.
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

// Load: reads the indexMap from a file with gob
func Load(filename string) (IndexMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var indexMap IndexMap
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&indexMap); err != nil {
		return nil, err
	}
	return indexMap, nil
}

// Lookup handles the entire lookup workflow: load index, search hash, print results
// 1.Load index from file
// 2.Perform lookup
// 3.Print results
func LookUp(input_file string, simHash string) error {
	indexmap, err := Load(input_file)
	if err != nil {
		return fmt.Errorf("Error loading index: %v\n", err)
	}

	entries, ok := indexmap[simHash]
	if !ok {
		return fmt.Errorf("No entries found for SimHash: %s\n", simHash)
	}
	LookUpOutput(simHash, entries)

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
