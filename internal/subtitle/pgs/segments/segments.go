package segments

import (
	"image/color"
	"time"
)

type SegmentType int

const (
	PDS SegmentType = iota
	ODS
	PCS
	WDS
	END
)

const INV SegmentType = -1

type Header struct {
	PTS         time.Duration
	SegmentType SegmentType
	SegmentSize int
}
type LastInSquenceFlag int

const (
	LastInSequence LastInSquenceFlag = iota
	FirstInSequence
	FirstAndLastInSequence
)

type ObjectDefinitionSegment struct {
	ObjectID            int
	ObjectVersionNumber int
	LastInSquenceFlag   LastInSquenceFlag
	Width               int
	Height              int
	ObjectData          []byte
}

type CompositionState int

const (
	EpochStart CompositionState = iota
	AcquisitionStart
	Normal
)

type PresentationCompositionSegment struct {
	Width              int
	Height             int
	CompositionState   CompositionState
	PaletteUpdateFlag  bool
	PaletteID          int
	CompositionObjects []CompositionObject
}

type CompositionObject struct {
	ObjectID                         int
	WindowID                         int
	ObjectCroppedFlag                bool
	ObjectHorizontalPosition         int
	ObjectVerticalPosition           int
	ObjectCroppingHorizontalPosition int
	ObjectCroppingVerticalPosition   int
	ObjectCroppingWidth              int
	ObjectCroppingHeight             int
}

type PaletteDefinitionSegment struct {
	PaletteID            int
	PaletteVersionNumber int
	PaletteEntries       map[int]color.Color
}

type Window struct {
	WindowID                 int
	WindowHorizontalPosition int
	WindowVerticalPosition   int
	WindowWidth              int
	WindowHeight             int
}

type WindowDefinitionSegment struct {
	Windows []Window
}
