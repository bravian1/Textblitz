package internals

import (
	"encoding/gob"
	"fmt"
	"os"
)

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

// save: writes the indexMap to  a file with csv
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
