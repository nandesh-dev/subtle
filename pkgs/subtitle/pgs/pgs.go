package pgs

import (
	"image"
	"time"
)

type Segment struct {
	start  time.Duration
	images []image.Image
}

func NewSegment() *Segment {
	return &Segment{}
}

func (s *Segment) Images() ([]image.Image, error) {
	return s.images, nil
}

func (s *Segment) AddImages(images []image.Image) error {
	s.images = append(s.images, images...)
	return nil
}

func (s *Segment) Start() time.Duration {
	return s.start
}

func (s *Segment) SetStart(start time.Duration) {
	s.start = start
}

type Stream struct {
	segments []Segment
}

func NewStream() *Stream {
	return &Stream{
		segments: make([]Segment, 0),
	}
}

func (s *Stream) AddSegment(segment Segment) {
	s.segments = append(s.segments, segment)
}

func (s *Stream) Segments() []Segment {
	return s.segments
}
