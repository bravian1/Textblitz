package main

import (
	"fmt"
	"github.com/bravian1/Textblitz/internals"
)

func main() {
	config, err := internals.ParseFlags()
	if err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	switch config.Command {
	case "index":
		fmt.Println("Performing indexing...")
		if err := internals.IndexFile(config.InputFile, config.ChunkSize, config.WorkerPool, config.OutputFile); err != nil {
			fmt.Printf("Error during indexing: %v\n", err)
			return
		}
		fmt.Printf("Successfully indexed %s\n", config.InputFile)
	case "lookup":
		fmt.Println("Performing lookup...")

		indexManager := internals.NewIndexManager()

		if err := indexManager.Load(config.InputFile); err != nil {
			fmt.Printf("Error loading index: %v\n", err)
			return
		}

		entries, err := indexManager.Lookup(config.SimHash)
		if err != nil {
			fmt.Printf("Error during lookup: %v\n", err)
			return
		}

		internals.LookUpOutput(config.SimHash, entries)
	default:
		fmt.Println("Invalid command. Use 'index' or 'lookup'.\n or --help for more information.")
	}
}
