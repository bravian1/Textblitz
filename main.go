package main

import (
	"fmt"
	"github.com/bravian1/Textblitz/internals"
	"runtime/pprof"
	"os"
)

func main() {
	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

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
		if err := internals.LookUp(config.SimHash, config.OutputFile); err != nil {
			fmt.Printf("Error during lookup: %v\n", err)
			return
		}
	default:
		fmt.Println("Invalid command. Use 'index' or 'lookup'.\n or --help for more information.")
	}
	f2, _ := os.Create("memory.pprof")
	pprof.WriteHeapProfile(f2)
	defer f2.Close()
}
