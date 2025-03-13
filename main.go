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
		fmt.Println("Performing indexing...\n")
		if err := internals.SystemIntegration(config.InputFile, config.ChunkSize, config.WorkerPool); err != nil {
			fmt.Printf("Error during indexing: %v\n", err)
			return
		}
		fmt.Printf("Successfully indexed %s\n", config.InputFile)
	case "lookup":
		fmt.Println("Performing lookup...\n")
	default:
		fmt.Println("Invalid command. Use 'index' or 'lookup'.\n or --help for more information.")
	}
}
