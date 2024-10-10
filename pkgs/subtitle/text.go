package subtitle

import "time"

type TextSubtitle struct {
	segments []TextSegment
}

type TextSegment struct {
	start time.Duration
	end   time.Duration
	text  string
}

func NewTextSubtitle() *TextSubtitle {
	return &TextSubtitle{
		segments: make([]TextSegment, 0),
	}
}

func (s *TextSubtitle) AddSegment(segment TextSegment) {
	s.segments = append(s.segments, segment)
}

func (s *TextSubtitle) Segments() []TextSegment {
	return s.segments
}

func NewTextSegment(start time.Duration, end time.Duration, text string) *TextSegment {
	return &TextSegment{
		start: start,
		end:   end,
		text:  text,
	}
}

func (s *TextSegment) Start() time.Duration {
	return s.start
}

func (s *TextSegment) End() time.Duration {
	return s.end
}

func (s *TextSegment) Text() string {
	return s.text
}
