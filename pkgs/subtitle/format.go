package subtitle

import "fmt"

type Format int

const (
	PGS Format = iota
	ASS
	SRT
)

func ParseFormat(f string) (Format, error) {
	switch f {
	case "ass":
		return ASS, nil
	case "pgs":
		return PGS, nil
	case "srt":
		return SRT, nil
	}

	return ASS, fmt.Errorf("Invalid format: %v", f)
}
