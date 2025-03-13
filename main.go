package main

import (
	"fmt"

	"github.com/bravian1/Textblitz/internals"
)

func main() {

	config := internals.ParseFlags()

	switch config.Command {
	case "index":
		fmt.Println("Performing indexing...\n")
	case "lookup":
		fmt.Println("Performing lookup...\n")
	default:
		fmt.Println("Invalid command. Use 'index' or 'lookup'.\n or --help for more information.")
	}
}
