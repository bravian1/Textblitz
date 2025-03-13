package internals

//index entry struct
type IndexEntry struct {
	SimHash uint64
	OriginalFile string
	Position int64 //offset
	AssociatedWords []string
}


type IndexManager interface{
	 Load(input_file string) //add to in memory map
	 Lookup(simhash uint64) ([]IndexEntry, error)
	 Add(entry IndexEntry) error //add index entry to map
	 Save(map[string]IndexEntry) error //save map to idx file
	
}