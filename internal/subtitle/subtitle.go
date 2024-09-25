package subtitle

import (
	"time"

	"golang.org/x/text/language"
)

type SubtitleFile struct {
	Path   string
	Format SubtitleFileFormat
}

type SubtitleFileFormat int

const (
	SRT SubtitleFileFormat = iota
	ASS
	SSA
	IDX
	SUB
	PGS
)

type SubtitleSegment struct {
	Start *time.Duration
	End   *time.Duration
	Text  string
	Style struct{}
}

type SubtitleStream struct {
	Langauge language.Tag
	Segments []SubtitleSegment
}

type Subtitle struct {
	Streams []SubtitleStream
}

type RawSubtitleStream struct {
	Index         int
	Format        string
	Language      language.Tag
	VideoFilePath string
}
