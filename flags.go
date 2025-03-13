package cmd

import "flag"

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
	//flags
	command := flag.String("c", "", "Command: 'index' to index a file, 'lookup to search a hash ")
	inputFile := flag.String("i", "", "Input file(text file for  index, .idx for  lookup)")
    chunkSize := flag.Int("s", 1024, "Chunk size in bytes (default 4096)")
	outputFile := flag.String("o", "", "Output index file (.idx) .Required for 'index' command")
	SimHash := flag.String("h", "", "Simhash value to search (required for 'lookup' command)")
	help := flag.Bool("help", false, "Display help message")

//parse flags
flag.Parse()

}

