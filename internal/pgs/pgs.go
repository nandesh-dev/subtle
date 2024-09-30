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

func (s *Segment) Images() []image.Image {
	return s.images
}

func (s *Segment) AddImages(img []image.Image) {
	s.images = append(s.images, img...)
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
