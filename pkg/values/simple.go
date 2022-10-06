package values

import "strings"

type Simple struct {
	Values []string
}

func (s *Simple) Set(val string) (err error) {
	switch {
	case strings.ContainsRune(val, ','):
		s.Values = append(s.Values, strings.Split(val, ",")...)
	default:
		s.Values = append(s.Values, val)
	}

	return
}

func (s *Simple) String() (rv string) {
	return strings.Join(s.Values, ",")
}
