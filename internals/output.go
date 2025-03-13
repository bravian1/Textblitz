package internals

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//header
	header := []string{"SimHash", "Original File", "Position", "Associated Words"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for simHash, entries := range indexmap {
		for _, entry := range entries {
			sizeStr := strconv.Itoa(entry.Size)
			words := strings.Join(entry.AssociatedWords, " ")
			record := []string{simHash, entry.OriginalFile, sizeStr, words}
			if err := writer.Write(record); err != nil {
				return err
			}
		}

	}
	return nil
}
