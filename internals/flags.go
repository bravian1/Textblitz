package internals

import (
	"flag"
	"fmt"
	"os"
)

// CLIflags holds the parsed command line arguments
type CLIFlags struct {
	Command    string //index/look up
	InputFile  string //path to .txt file (for  index) or .idx file (for look up)
	ChunkSize  int    //chunk size (bytes)
	OutputFile string //path to output.idx file
	SimHash    string //simhash value to search
	WorkerPool int    //number of worker goroutines
	Threshold int // distance for fuzzy lookup
	
}

// Parseflags parses command line arguments and returns a CLIFlags struct
func ParseFlags() (CLIFlags, error) {
	config := CLIFlags{}
	flagSet := flag.NewFlagSet("textblitz", flag.ExitOnError)

	//flags
	flagSet.StringVar(&config.Command, "c", "", "Command: 'index' to index a file, 'lookup to search a hash ")
	flagSet.StringVar(&config.InputFile, "i", "", "Input file(text file for  index, .idx for  lookup)")
	flagSet.IntVar(&config.ChunkSize, "s", 4096, "Chunk size in bytes (default 4096)")
	flagSet.StringVar(&config.OutputFile, "o", "", "Output index file (.idx) .Required for 'index' command")
	flagSet.StringVar(&config.SimHash, "h", "", "Simhash value to search (required for 'lookup' command)")
	flagSet.IntVar(&config.WorkerPool, "w", 4, "Number of worker goroutines (default 4)")
	flagSet.IntVar(&config.Threshold, "t", 0, "Distance for fuzzy lookup (default 0)")
	help := flagSet.Bool("help", false, "Display help message")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return config, err
	}

	if *help {
		PrintHelp()
		return config, fmt.Errorf("help message displayed")
	}

	//validate flags
	if config.Command == "" {
		return config, fmt.Errorf("Error: Missing Command  (-c 'index' or 'lookup'). Use --help for details.")
	}

	if config.Command == "index" && (config.InputFile == "" || config.OutputFile == "") {
		return config, fmt.Errorf("Error: Input file  (-i <input_file.txt> )or OutputFile (-o <index.idx>)  are required for indexing. Use --help for details.")
	}

	if config.Command == "lookup" && (config.InputFile == "" || config.SimHash == "") {
		return config, fmt.Errorf("Error:Input file  (-i <index_file.idx>) or SimHash (-h <simhash_value>)  are required for lookup. Use --help for details.")
	}

	return config, nil
}

// print help message
func PrintHelp() {
	fmt.Println(`TextIndex CLI - Fast & Scalable Text Indexer
--------------------------------------------
A command-line tool for indexing large text files and performing fast lookups using SimHash.

Usage:
  textindex -c index -i <input_file> -s <chunk_size> -o <index_file> [-w <workers>]
  textindex -c lookup -i <index_file> -h <simhash_value>

Commands:
  -c index   : Index a file by splitting it into chunks, computing SimHash, and saving the index.
  -c lookup  : Find a chunk in the indexed file based on its SimHash.

Arguments:
  -i <file>      : Input file (text file for indexing, .idx file for lookup).
  -s <size>      : Chunk size in bytes (default: 4096).
  -o <file>      : Output index file (required for indexing).
  -h <simhash>   : SimHash value to search for (required for lookup).
  -w <workers>   : Number of workers (Goroutines) for parallel indexing (default: 4).
  -t <threshold> : Distance for fuzzy lookup (default 0).
  --help         : Display this help message.

Example Usage:
  # Index a file with 4KB chunks using 4 workers
  textindex -c index -i large_text.txt -s 4096 -o index.idx -w 4

  # Lookup a SimHash value in an index file with a threshold of 2
  textindex -c lookup -i index.idx -h 3e4f1b2c98a6 -t 2

Error Handling:
  - "File not found"  : Ensure the input file exists.
  - "Invalid chunk size" : Use a valid numeric chunk size (e.g., 1024, 4096).
  - "SimHash not found" : Ensure the index file was generated before lookup.

For more details, refer to the README.md`)
}
