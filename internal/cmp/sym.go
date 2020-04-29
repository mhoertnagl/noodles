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

// NewClosureSymTable create a closure symbol table. This table has the same
// parent as sym itself. It is NOT a child of sym. The closure table contains
// the closure params beginning with index 0 and then all params that are in
// sym shifted in index by the number of closure params.
func NewClosureSymTable(s *SymTable, eps []*SymbolNode) *SymTable {
	// The closure symbol table is NOT a child of s but of the parent of s.
	cs := &SymTable{
		parent:  s.parent,
		entries: make(symMap),
	}
	// Add all the new closure parameters beginning at index 0.
	for idx, ps := range eps {
		cs.entries[ps.Name] = &SymEntry{
			idx: idx,
		}
	}
	// Shift all local parameter indexes by the number of closure parameters.
	cargs := len(eps)
	for name, n := range s.entries {
		cs.entries[name] = &SymEntry{
			idx: cargs + n.idx,
		}
	}
	return cs
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

// func (s *SymTable) AddClosureParams(ns []*SymbolNode) {
// 	cargs := len(ns)
// 	// Shift all local parameter indexes by the number of closure parameters.
// 	for _, n := range s.entries {
// 		n.idx += cargs
// 	}
// 	// Add all the new closure parameters beginning at index 0.
// 	for idx, n := range ns {
// 		s.entries[n.Name] = &SymEntry{
// 			idx: idx,
// 		}
// 	}
// }

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
