package internals

import "fmt"

//Formats the outputs for index and lookup
func LookUpOutput(entries []IndexEntry) {
	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return
	}

	fmt.Println("\nLookup Complete!")
	fmt.Println("------------------------------------")

	for _, entry := range entries {
		fmt.Printf("| SimHash       : %s\n", entry.SimHash)
		fmt.Printf("| Original File : %s\n", entry.OriginalFile)
		fmt.Printf("| Position      : Byte %d\n", entry.Position)
		fmt.Printf("| Associated Words : \"%s\"\n", entry.AssociatedWords)
		fmt.Println("------------------------------------------------")
	}

	fmt.Println()
}
