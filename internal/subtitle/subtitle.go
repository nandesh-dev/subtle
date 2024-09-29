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

func NewStream(lang language.Tag) Stream {
	return Stream{
		Langauge: lang,
		Segments: make([]Segment, 0),
	}
}

func (s *Stream) AddSegment(seg Segment) {
	if len(s.Segments) > 0 {
		previousSegment := s.Segments[len(s.Segments)-1]
		if seg.Start == nil && previousSegment.End != nil {
			seg.Start = &*previousSegment.End
		} else if seg.Start != nil && previousSegment.End == nil {
			previousSegment.End = &*seg.Start
		}
	}

	s.Segments = append(s.Segments, seg)
}

type Subtitle struct {
	Streams []Stream
}

type RawStream struct {
	Index         int
	Format        Format
	Language      language.Tag
	VideoFilePath string
}
