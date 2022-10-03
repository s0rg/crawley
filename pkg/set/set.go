package set

type Set[T comparable] map[T]stub

// Add add value to set, replacing previous instances.
func (s Set[T]) Add(v T) {
	s[v] = stub{}
}

// Has checks if value is already present in set.
func (s Set[T]) Has(v T) (ok bool) {
	_, ok = s[v]

	return
}

// List returns set as slice of unique strings.
func (s Set[T]) List() (rv []T) {
	rv = make([]T, 0, len(s))

	for k := range s {
		rv = append(rv, k)
	}

	return rv
}

// Load populates set with given values.
func (s Set[T]) Load(vals ...T) Set[T] {
	for _, v := range vals {
		s.Add(v)
	}

	return s
}
