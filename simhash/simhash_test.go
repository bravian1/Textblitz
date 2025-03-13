package simhash

import (
	"testing"
)

func TestWordFeatureSet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Feature
	}{
		{
			name:  "Simple sentence",
			input: "Hello world!",
			expected: []Feature{
				{Text: "hello", Weight: 1},
				{Text: "world", Weight: 1},
			},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []Feature{},
		},
		{
			name:  "Multiple spaces and punctuation",
			input: "Hello,   world! How are you?",
			expected: []Feature{
				{Text: "hello", Weight: 1},
				{Text: "world", Weight: 1},
				{Text: "how", Weight: 1},
				{Text: "are", Weight: 1},
				{Text: "you", Weight: 1},
			},
		},
	}

	fs := NewWordFeatureSet()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			features := fs.Features(tt.input)
			if len(features) != len(tt.expected) {
				t.Errorf("expected %d features, got %d", len(tt.expected), len(features))
				return
			}
			for i, feature := range features {
				if feature.Text != tt.expected[i].Text {
					t.Errorf("feature %d: expected text %q, got %q", i, tt.expected[i].Text, feature.Text)
				}
				if feature.Weight != tt.expected[i].Weight {
					t.Errorf("feature %d: expected weight %d, got %d", i, tt.expected[i].Weight, feature.Weight)
				}
			}
		})
	}
}

func TestNgramFeatureSet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		step     int
		expected []Feature
	}{
		{
			name:  "Simple trigram",
			input: "Hello",
			n:     3,
			step:  1,
			expected: []Feature{
				{Text: "hel", Weight: 1},
				{Text: "ell", Weight: 1},
				{Text: "llo", Weight: 1},
			},
		},
		{
			name:     "Empty string",
			input:    "",
			n:        3,
			step:     1,
			expected: []Feature{},
		},
		{
			name:  "Step size 2",
			input: "Hello",
			n:     3,
			step:  2,
			expected: []Feature{
				{Text: "hel", Weight: 1},
				{Text: "llo", Weight: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewNgramFeatureSet(tt.n, tt.step)
			features := fs.Features(tt.input)
			if len(features) != len(tt.expected) {
				t.Errorf("expected %d features, got %d", len(tt.expected), len(features))
				return
			}
			for i, feature := range features {
				if feature.Text != tt.expected[i].Text {
					t.Errorf("feature %d: expected text %q, got %q", i, tt.expected[i].Text, feature.Text)
				}
				if feature.Weight != tt.expected[i].Weight {
					t.Errorf("feature %d: expected weight %d, got %d", i, tt.expected[i].Weight, feature.Weight)
				}
			}
		})
	}
}

func TestSimHash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fs       FeatureSet
		expected uint64
	}{
		{
			name:  "Word-based hash",
			input: "Hello world",
			fs:    NewWordFeatureSet(),
		},
		{
			name:  "Ngram-based hash",
			input: "Hello world",
			fs:    NewNgramFeatureSet(3, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewSimHashGenerator(tt.fs)
			hash := gen.Hash(tt.input)
			if hash == 0 {
				t.Error("expected non-zero hash")
			}
		})
	}
}

func TestSimHashSimilarity(t *testing.T) {
	gen := NewSimHashGenerator(NewWordFeatureSet())

	// Test similar texts
	text1 := "The quick brown fox jumps over the lazy dog"
	text2 := "The quick brown fox jumps over the lazy dog"
	hash1 := gen.Hash(text1)
	hash2 := gen.Hash(text2)

	if hash1 != hash2 {
		t.Error("identical texts should produce identical hashes")
	}

	// Test different texts
	text3 := "A completely different text that should produce a different hash"
	hash3 := gen.Hash(text3)

	if hash1 == hash3 {
		t.Error("different texts should produce different hashes")
	}
}
