package internals

type IndexEntry struct {
	OriginalFile    string
	Size            int
	Position        int
	AssociatedWords []string
}

type IndexMap map[string][]IndexEntry

type IndexManager interface {
	Load(inputFile string) error

	Lookup(simhash string) ([]IndexEntry, error)

	Add(simhash string, entry IndexEntry) error

	Save(outputFile string) error
}

type indexManager struct {
	index IndexMap
}
