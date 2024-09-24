package subtitle

import (
	"image"
	"time"
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

type ImageSubtitleSegment struct {
	Start  time.Duration
	Images []image.Image
}

type Subtitle struct {
	ImageStream []ImageSubtitleSegment
}
