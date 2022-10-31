package set

type (
	stub struct{}

	// Set represents hashset for comparable types.
	Set[T comparable] map[T]stub
)

// Add add value to set, replacing previous instances.
func (s Set[T]) Add(v T) {
	s[v] = stub{}
}

// TryAdd takes attempt to add value to set, returns false if value already exists.
func (s Set[T]) TryAdd(v T) (ok bool) {
	if s.Has(v) {
		return false
	}

	s.Add(v)

	return true
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
