package ass

import (
	"fmt"
	"strings"
)

type reader struct {
	Lines    []string
	Position int
}

func NewReader(data []byte) *reader {
	str := string(data)

	lines := strings.Split(str, "\r\n")

	if len(lines) == 1 {
		lines = strings.Split(str, "\n")
	}

	return &reader{
		Lines:    lines,
		Position: 0,
	}
}

func (r *reader) Advance() (string, error) {
	if r.ReachedEnd() {
		return "", fmt.Errorf("Already reached the end")
	}

	r.Position++

	return r.Lines[r.Position-1], nil
}

func (r *reader) ReachedEnd() bool {
	return r.Position >= len(r.Lines)
}
