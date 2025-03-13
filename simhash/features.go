package simhash

// Feature represents a piece of text with its importance weight.
type Feature struct {
	Text   string
	Weight int
}

// FeatureSet defines how to break down text into meaningful features.
type FeatureSet interface {
	Features(text string) []Feature
}
