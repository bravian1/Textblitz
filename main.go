package main

import (
	"fmt"

	cmd "github.com/bravian1/Textblitz/internals"
)

func main() {

	//command line flags
	config := cmd.ParseFlags()

	//handle commands using the struct
	switch config.Command {
	case "index":
		fmt.Println("Performing indexing...\n")
	case "lookup":
		fmt.Println("Performing lookup...\n")
	default:
		fmt.Println("Invalid command. Use 'index' or 'lookup'.\n or --help for more information.")
	}
}
