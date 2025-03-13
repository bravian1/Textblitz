package simhash

type SimHashGen struct {
	FeatureSet FeatureSet
}

func NewSimHashGenerator(fs FeatureSet) *SimHashGen {
	return &SimHashGen{FeatureSet: fs}
}
