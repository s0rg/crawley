package values

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
)

const fileMarker = '@' // curl-compatible

type List struct {
	values []string
}

func (l *List) String() (rv string) {
	return
}

func (l *List) Set(val string) (err error) {
	l.values = append(l.values, val)

	return
}

func (l *List) Load(
	target fs.FS,
) (rv []string, err error) {
	rv = make([]string, 0, len(l.values))

	var vals []string

	for _, v := range l.values {
		if v[0] == fileMarker {
			if vals, err = l.loadFile(target, v[1:]); err != nil {
				return
			}

			rv = append(rv, vals...)
		} else {
			rv = append(rv, v)
		}
	}

	return
}

func (l *List) loadFile(
	target fs.FS,
	name string,
) (rv []string, err error) {
	var body []byte

	if body, err = fs.ReadFile(target, name); err != nil {
		err = fmt.Errorf("read: %w", err)

		return
	}

	s := bufio.NewScanner(bytes.NewReader(body))

	for s.Scan() {
		rv = append(rv, s.Text())
	}

	return rv, nil
}
