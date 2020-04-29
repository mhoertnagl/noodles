package cmp

type SymEntry struct {
	idx int
}

type symMap map[string]*SymEntry

type SymTable struct {
	// Name    string
	parent  *SymTable
	entries symMap
}

func NewSymTable() *SymTable {
	return &SymTable{
		parent:  nil,
		entries: make(symMap),
	}
}

func (s *SymTable) NewSymTable() *SymTable {
	return &SymTable{
		parent:  s,
		entries: make(symMap),
	}
}

func (s *SymTable) Size() int {
	return len(s.entries)
}

func (s *SymTable) AddVar(ns ...string) {
	s.Add(ns)
}

func (s *SymTable) Add(ns []string) {
	sz := s.Size()
	for idx, n := range ns {
		s.entries[n] = &SymEntry{
			idx: sz + idx,
		}
	}
}

func (s *SymTable) RemoveVar(ns ...string) {
	s.Remove(ns)
}

func (s *SymTable) Remove(ns []string) {
	for _, n := range ns {
		delete(s.entries, n)
	}
}

// func (s *SymTable) Find(n string) (*SymEntry, bool) {
// 	for c := s; c != nil; c = c.parent {
// 		if e, ok := c.entries[n]; ok {
// 			return e, true
// 		}
// 	}
// 	return nil, false
// }

func (s *SymTable) IndexOf(n string) (int, bool) {
	dfp := 0
	for c := s; c != nil; c = c.parent {
		if e, ok := c.entries[n]; ok {
			return dfp + e.idx, true
		}
		// Subtract the FP and the RP cell as well as the number of arguments
		// of the current frame.
		dfp -= 2 + c.Size()
	}
	return 0, false
}
