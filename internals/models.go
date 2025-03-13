package internals

//index entry struct
type IndexEntry struct {
	OriginalFile string
	Size int
	Position int //offset
	AssociatedWords []string
}

type IndexMap map[string][]IndexEntry 

// type IndexManager interface{
// 	 Load(input_file string) //add to in memory map
// 	 Lookup(simhash uint64) ([]IndexEntry, error)
// 	 Add(entry IndexEntry) error //add index entry to map
// 	 Save(map[uint64][]IndexEntry) error //save map to idx file
	
// }

