package simhash

import "hash/fnv"

type SimHashGen struct {
	FeatureSet FeatureSet
}

func NewSimHashGenerator(fs FeatureSet) *SimHashGen {
	return &SimHashGen{FeatureSet: fs}
}

// Hash takes the chunk and turns it into a 64-bit SimHash number.
// Divide the chunk to features
// Hash each feature into a 64-bit number using FNV-1a .
// For each of the 64 bit positions:
//   - If the feature’s hash has a 1 in that spot, it adds the feature’s weight.
//   - If it’s a 0, it subtracts the weight.
//
// At the end, the SimHash has a 1 in any bit position where the total weight is positive.
func (sg *SimHashGen) Hash(text string) uint64 {
	features := sg.FeatureSet.Features(text)

	bitCounts := make([]int, 64)
	hasher := fnv.New64a()

	for _, feature := range features {
		hasher.Reset()
		hasher.Write([]byte(feature.Text))
		featureHash := hasher.Sum64()

		for i := range 64 {
			if (featureHash & (1 << i)) != 0 {
				bitCounts[i] += feature.Weight
			} else {
				bitCounts[i] -= feature.Weight
			}
		}
	}

	var simhash uint64
	for i := range 64 {
		if bitCounts[i] > 0 {
			simhash |= (1 << i)
		}
	}

	return simhash
}
