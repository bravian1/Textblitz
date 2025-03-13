package internals

import (
	"os"
	"testing"
)

// Helper function to reset os.Args and avoid conflicts between tests
func resetArgs(args []string) {
	os.Args = append([]string{"textindex"}, args...)
}

// Test `ParseFlags()` for a valid index command
func TestParseFlags_IndexCommand(t *testing.T) {
	resetArgs([]string{"-c", "index", "-i", "sample.txt", "-s", "4096", "-o", "index.idx", "-w", "8"})

	config := ParseFlags()

	if config.Command != "index" {
		t.Errorf("Expected command 'index', got %s", config.Command)
	}
	if config.InputFile != "sample.txt" {
		t.Errorf("Expected input file 'sample.txt', got %s", config.InputFile)
	}
	if config.ChunkSize != 4096 {
		t.Errorf("Expected chunk size 4096, got %d", config.ChunkSize)
	}
	if config.OutputFile != "index.idx" {
		t.Errorf("Expected output file 'index.idx', got %s", config.OutputFile)
	}
	if config.WorkerPool != 8 {
		t.Errorf("Expected worker pool 8, got %d", config.WorkerPool)
	}
}

// Test `ParseFlags()` for a valid lookup command
func TestParseFlags_LookupCommand(t *testing.T) {
	resetArgs([]string{"-c", "lookup", "-i", "index.idx", "-h", "3e4f1b2c98a6"})

	config := ParseFlags()

	if config.Command != "lookup" {
		t.Errorf("Expected command 'lookup', got %s", config.Command)
	}
	if config.InputFile != "index.idx" {
		t.Errorf("Expected input file 'index.idx', got %s", config.InputFile)
	}
	if config.SimHash != "3e4f1b2c98a6" {
		t.Errorf("Expected SimHash '3e4f1b2c98a6', got %s", config.SimHash)
	}
}

// Test missing command flag
func TestParseFlags_MissingCommand(t *testing.T) {
	resetArgs([]string{"-i", "sample.txt", "-s", "4096", "-o", "index.idx"})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected error for missing command, but function did not exit")
		}
	}()

	ParseFlags()
}

// Test missing required arguments for index command
func TestParseFlags_MissingIndexArgs(t *testing.T) {
	resetArgs([]string{"-c", "index", "-i", "sample.txt"})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected error for missing output file in index command, but function did not exit")
		}
	}()

	ParseFlags()
}

// Test missing required arguments for lookup command
func TestParseFlags_MissingLookupArgs(t *testing.T) {
	resetArgs([]string{"-c", "lookup", "-i", "index.idx"})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected error for missing SimHash in lookup command, but function did not exit")
		}
	}()

	ParseFlags()
}

// Test `--help` flag
func TestParseFlags_HelpFlag(t *testing.T) {
	resetArgs([]string{"--help"})

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Help flag should exit gracefully, but function panicked")
		}
	}()

	ParseFlags()
}
