package values

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
)

const fileMarker = '@' // curl-compatible

type Smart struct {
	values []string
}

func (s *Smart) Set(val string) (err error) {
	s.values = append(s.values, val)

	return
}

func (s *Smart) Load(
	target fs.FS,
) (rv []string, err error) {
	rv = make([]string, 0, len(s.values))

	var vals []string

	for _, v := range s.values {
		if v[0] == fileMarker {
			if vals, err = loadFile(target, v[1:]); err != nil {
				return
			}

			rv = append(rv, vals...)
		} else {
			rv = append(rv, v)
		}
	}

	return
}

func (s *Smart) String() (rv string) {
	return
}

func loadFile(
	target fs.FS,
	name string,
) (rv []string, err error) {
	var body []byte

	if body, err = fs.ReadFile(target, name); err != nil {
		err = fmt.Errorf("read: %w", err)

		return
	}

	sc := bufio.NewScanner(bytes.NewReader(body))

	for sc.Scan() {
		rv = append(rv, sc.Text())
	}

	return rv, nil
}
