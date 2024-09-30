package ass

import "time"

type Segment struct {
	start time.Duration
	end   time.Duration
	text  string
}

func NewSegment() *Segment {
	return &Segment{}
}

func (s *Segment) Text() string {
	return s.text
}

func (s *Segment) SetText(text string) {
	s.text = text
}

func (s *Segment) Start() time.Duration {
	return s.start
}

func (s *Segment) SetStart(start time.Duration) {
	s.start = start
}

func (s *Segment) End() time.Duration {
	return s.end
}

func (s *Segment) SetEnd(end time.Duration) {
	s.end = end
}

type Stream struct {
	segments []Segment
}

func NewStream() *Stream {
	return &Stream{
		segments: make([]Segment, 0),
	}
}

func (s *Stream) Segments() []Segment {
	return s.segments
}
