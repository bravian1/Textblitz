package internals

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"
)

// TestSave verifies that the Save function correctly writes index data to a CSV file.
func TestSave(t *testing.T) {
	// Sample IndexMap data
	indexmap := IndexMap{
		"3e4f1b2c": {
			{OriginalFile: "large_text.txt", Size: 4096, AssociatedWords: []string{"Once", "upon", "a", "time"}},
			{OriginalFile: "another_text.txt", Size: 4096, AssociatedWords: []string{"The", "story", "continues"}},
		},
		"a7c9d4f8": {
			{OriginalFile: "large_text.txt", Size: 4096, AssociatedWords: []string{"In", "a", "faraway", "land"}},
		},
	}

	// Create a temporary file for testing
	tempFile := "test_index.idx"
	defer os.Remove(tempFile) // Clean up after test

	// Call Save function
	err := Save(tempFile, indexmap)
	if err != nil {
		t.Fatalf("Save function failed: %v", err)
	}

	// Open the file for reading
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	// Validate header row
	expectedHeader := []string{"SimHash", "Original File", "Position", "Associated Words"}
	if !reflect.DeepEqual(records[0], expectedHeader) {
		t.Errorf("Header mismatch. Expected: %v, Got: %v", expectedHeader, records[0])
	}

	// Validate first data row
	expectedFirstRow := []string{
		"3e4f1b2c", "large_text.txt", "4096", "Once upon a time",
	}
	if !reflect.DeepEqual(records[1], expectedFirstRow) {
		t.Errorf("First row mismatch. Expected: %v, Got: %v", expectedFirstRow, records[1])
	}

	// Validate second data row (multiple entries for the same SimHash)
	expectedSecondRow := []string{
		"3e4f1b2c", "another_text.txt", "4096", "The story continues",
	}
	if !reflect.DeepEqual(records[2], expectedSecondRow) {
		t.Errorf("Second row mismatch. Expected: %v, Got: %v", expectedSecondRow, records[2])
	}

	// Validate third row
	expectedThirdRow := []string{
		"a7c9d4f8", "large_text.txt", "4096", "In a faraway land",
	}
	if !reflect.DeepEqual(records[3], expectedThirdRow) {
		t.Errorf("Third row mismatch. Expected: %v, Got: %v", expectedThirdRow, records[3])
	}
}
