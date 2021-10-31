package set

// String holds set of string values.
type String map[string]stub

// Add add value to set, replacing previous instances.
func (ss String) Add(v string) {
	ss[v] = stub{}
}

// List returns set as slice of unique strings.
func (ss String) List() (rv []string) {
	rv = make([]string, 0, len(ss))

	for k := range ss {
		rv = append(rv, k)
	}

	return rv
}
