package simhash

import (
	"strings"
	"unicode"
)

// WordFeatureSet breaks text down into individual words.
type WordFeatureSet struct {
	Normalize bool
}

// NewWordFeatureSet creates a new word-based feature extractor.
// By default, we normalize everything to lowercase
func NewWordFeatureSet() *WordFeatureSet {
	return &WordFeatureSet{Normalize: true}
}

// Features takes a chunk of text and breaks it into individual words.
// It works like this:
// - First, make everything lowercase if normalization is on
// - Next, split the text by anything that's not a letter or number
// - Finally, create a Feature for each word with a weight of 1
func (w *WordFeatureSet) Features(text string) []Feature {
	if w.Normalize {
		text = strings.ToLower(text)
	}
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	features := make([]Feature, 0, len(words))
	for _, word := range words {
		if len(word) > 0 {
			features = append(features, Feature{Text: word, Weight: 1})
		}
	}
	return features
}
