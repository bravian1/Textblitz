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
	command := flag.String("c", "", "Command: 'index' to index a file, ")
}

