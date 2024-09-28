package subtitle

import (
	"time"

	"golang.org/x/text/language"
)

type Segment struct {
	Start *time.Duration
	End   *time.Duration
	Text  string
	Style struct{}
}

type Stream struct {
	Langauge language.Tag
	Segments []Segment
}

type Subtitle struct {
	Streams []Stream
}

type RawStream struct {
	Index         int
	Format        string
	Language      language.Tag
	VideoFilePath string
}
