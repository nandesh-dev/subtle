package subtitle

import (
	"image"
	"time"
)

type ImageSubtitle struct {
	segments []ImageSegment
}

type ImageSegment struct {
	start time.Duration
	end   time.Duration
	image image.Image
}

func NewImageSubtitle() *ImageSubtitle {
	return &ImageSubtitle{
		segments: make([]ImageSegment, 0),
	}
}

func (s *ImageSubtitle) AddSegment(segment ImageSegment) {
	s.segments = append(s.segments, segment)
}

func (s *ImageSubtitle) Segments() []ImageSegment {
	return s.segments
}

func NewImageSegment(start time.Duration, end time.Duration, image image.Image) *ImageSegment {
	return &ImageSegment{
		start: start,
		end:   end,
		image: image,
	}
}

func (s *ImageSegment) Start() time.Duration {
	return s.start
}

func (s *ImageSegment) End() time.Duration {
	return s.end
}

func (s *ImageSegment) Image() image.Image {
	return s.image
}
