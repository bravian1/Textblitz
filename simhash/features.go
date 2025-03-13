package simhash

type Feature struct {
    Text   string
    Weight int
}

type FeatureSet interface {
    Features(text string) []Feature
}


