package simhash

import (
	"strings"
    "unicode"
)

type WordFeatureSet struct {
	Normalize bool
}

func NewWordFeatureSet() *WordFeatureSet {
	return &WordFeatureSet{Normalize: true}
}

func (w *WordFeatureSet) Features(text string) []Feature {
	if w.Normalize {
		text = strings.ToLower(text)
	}
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	
	features := make([]Feature, len(words))
	for _ , word := range words {
		if len(word) > 0 {
			features = append(features, Feature{Text: word, Weight: 1})
		}
	}
	return features	
}