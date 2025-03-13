package cmd

import (
	"flag"
	"fmt"
	"log"
)

//CLIflags holds the parsed command line arguments
type CLIFlags struct {
	Command string //index/look up
	InputFile string //path to .txt file (for  index) or .idx file (for look up)
	ChunkSize int //chunk size (bytes)
	OutputFile string //path to output.idx file
	SimHash string //simhash value to search	
}

//Parseflags parses command line arguments and returns a CLIFlags struct
func ParseFlags() CLIFlags {
	config := CLIFlags{}

	//flags
	flag.StringVar(&config.Command,"c", "", "Command: 'index' to index a file, 'lookup to search a hash ")
	flag.StringVar(&config.InputFile,"i", "", "Input file(text file for  index, .idx for  lookup)")
    flag.IntVar(&config.ChunkSize,"s", 4096, "Chunk size in bytes (default 4096)")
	flag.StringVar(&config.OutputFile,"o", "", "Output index file (.idx) .Required for 'index' command")
	flag.StringVar(&config.SimHash,"h", "", "Simhash value to search (required for 'lookup' command)")
	help := flag.Bool("help", false, "Display help message")

//parse flags
flag.Parse()

//validate flags
if config.Command == "" {
	log.Fatalf("Error: Missing Command  (-c 'index' or 'lookup'). Use --help for details.")
}

if config.Command == "index" && (config.InputFile == "" || config.OutputFile == "") {
	log.Fatalf("Error: Input file  (-i <input_file.txt> )or OutputFile (-o <index.idx>)  are required for indexing. Use --help for details.")
}

if config.Command == "lookup" && (config.InputFile == "" || config.SimHash == "") {
	log.Fatalf("Error:Input file  (-i <index_file.idx>) or SimHash (-h <simhash_value>)  are required for lookup. Use --help for details.")
}

}