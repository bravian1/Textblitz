package simhash

import (
	"strings"	
)

type NgramFeatureSet struct {
	N int
	Step int
	Normalize bool
}

func NewNgramFeatureSet(n int, step int) *NgramFeatureSet {
	if n <= 0 {
		n = 3
	}

	if step <= 0 {
		step = 1
	}

	return &NgramFeatureSet{N: n, Step: step, Normalize: true}
}

func (ng *NgramFeatureSet) Features(text string) []Feature {
	if ng.Normalize {
		text = strings.ToLower(text)
	}

	if len(text) < ng.N {
		return []Feature{}
	}

	features := make([]Feature, 0, (len(text) - ng.N) / ng.Step + 1)
	for i := 0; i <= len(text) - ng.N; i += ng.Step {
		ngram := text[i:i+ng.N]
		features = append(features, Feature{Text: ngram, Weight: 1})
	}
	return features
	
}