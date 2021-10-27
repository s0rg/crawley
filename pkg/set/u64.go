package set

// U64 holds set of uint64 values.
type U64 map[uint64]stub

// Add try add value to set, returns true on success.
func (s U64) Add(v uint64) (ok bool) {
	if _, ok = s[v]; ok {
		return false
	}

	s[v] = stub{}

	return true
}
