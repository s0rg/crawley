package set

type U64 map[uint64]stub

func (s U64) Add(v uint64) (ok bool) {
	if _, ok = s[v]; ok {
		return false
	}

	s[v] = stub{}

	return true
}
