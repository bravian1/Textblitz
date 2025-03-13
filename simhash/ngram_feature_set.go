package simhash

import (
	"strings"
)

// NgramFeatureSet breaks text into overlapping chunks of n characters.
// It's like looking at text through a sliding window - we see small
// pieces at a time, sliding along to capture all the local patterns.
type NgramFeatureSet struct {
	// N is the size of each n-gram chunk
	N int
	// Step is how far to move the window each time
	Step int
	// Normalize determines if we convert everything to lowercase first
	Normalize bool
}

// NewNgramFeatureSet creates a new n-gram feature extractor.
// default values are n=3 and step=1
// step is the number of characters to shift the window each time
func NewNgramFeatureSet(n int, step int) *NgramFeatureSet {
	if n <= 0 {
		n = 3
	}

	if step <= 0 {
		step = 1
	}

	return &NgramFeatureSet{N: n, Step: step, Normalize: true}
}

// Features slices the text into overlapping n-grams.
// Make everything lowercase if normalization is on
// - Check if the text is long enough (must be at least n characters)
// - Then slide our window of size N across the text, stepping by Step each time
// - Each window position creates one feature with weight 1
func (ng *NgramFeatureSet) Features(text string) []Feature {
	if ng.Normalize {
		text = strings.ToLower(text)
	}

	if len(text) < ng.N {
		return []Feature{}
	}

	features := make([]Feature, 0, (len(text)-ng.N)/ng.Step+1)
	for i := 0; i <= len(text)-ng.N; i += ng.Step {
		ngram := text[i : i+ng.N]
		features = append(features, Feature{Text: ngram, Weight: 1})
	}
	return features
}
