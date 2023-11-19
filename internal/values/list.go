package values

import "strings"

type List struct {
	Values []string
}

func (s *List) Set(val string) (err error) {
	switch {
	case strings.ContainsRune(val, ','):
		s.Values = append(s.Values, strings.Split(val, ",")...)
	default:
		s.Values = append(s.Values, val)
	}

	return
}

func (s *List) String() (rv string) {
	return strings.Join(s.Values, ",")
}
