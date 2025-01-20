package subtitle

import (
	"image"
	"time"

	"golang.org/x/text/language"
)

type Subtitle struct {
	Metadata Metadata
	Cues     []Cue
}

type Writer interface {
	Write(*Subtitle) error
}

type Parser interface {
	Parse([]byte) (*Subtitle, error)
}

func (sub *Subtitle) Write(writer Writer) error {
	return writer.Write(sub)
}

type Metadata struct {
	Language language.Tag
}

type Cue struct {
	Timestamp      CueTimestamp
	Content        []CueContentSegment
	OriginalImages []image.Image
}

type CueTimestamp struct {
	Start time.Duration
	End   time.Duration
}

type CueContentSegment struct {
	Text  string
	Style CueContentSegmentStyle
}

type CueContentSegmentStyle struct{}
