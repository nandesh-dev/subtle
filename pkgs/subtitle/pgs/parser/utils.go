package parser

import (
	"fmt"
	"image/color"
	"time"
)

type ObjectDefinitionSequenceFlag int

const (
	LastInObjectDefinitionSequence ObjectDefinitionSequenceFlag = iota
	FirstInObjectDefinitionSequence
	FirstAndLastInObjectDefinitionSequence
	IntermediateInObjectDefinitionSequence
)

type ObjectDefinitionSegment struct {
	objectId            int
	objectVersionNumber int
	sequenceFlag        ObjectDefinitionSequenceFlag
	width               int
	height              int
	objectData          []byte
}

type PresentationCompositionState int

const (
	EpochStartPresentationCompositionState PresentationCompositionState = iota
	AcquisitionStartPresentationCompositionState
	NormalPresentationCompositionState
)

type PresentationCompositionObject struct {
	objectId                         int
	windowId                         int
	objectCroppedFlag                bool
	objectHorizontalPosition         int
	objectVerticalPosition           int
	objectCroppingHorizontalPosition int
	objectCroppingVerticalPosition   int
	objectCroppingWidth              int
	objectCroppingHeight             int
}

type PresentationCompositionSegment struct {
	width             int
	height            int
	state             PresentationCompositionState
	paletteUpdateFlag bool
	paletteId         int
	objects           []PresentationCompositionObject
}

type PaletteDefinitionSegment struct {
	paletteId            int
	paletteVersionNumber int
	paletteEntries       map[int]color.Color
}

type Window struct {
	id                 int
	horizontalPosition int
	verticalPosition   int
	width              int
	height             int
}

type WindowDefinitionSegment struct {
	windows []Window
}

type SegmentType int

const (
	PDSSegment SegmentType = iota
	ODSSegment
	PCSSegment
	WDSSegment
	ENDSegment
)

type Header struct {
	pts         time.Duration
	segmentType SegmentType
	segmentSize int
}

type DisplaySet struct {
	header                         Header
	presentationCompositionSegment PresentationCompositionSegment
	windowDefinitions              map[int]*Window
	paletteDefinitionSegments      map[int]*PaletteDefinitionSegment
	objectDefinitionSegments       map[int]*ObjectDefinitionSegment
}

func NewDisplaySet() *DisplaySet {
	return &DisplaySet{
		paletteDefinitionSegments: make(map[int]*PaletteDefinitionSegment),
		windowDefinitions:         make(map[int]*Window),
		objectDefinitionSegments:  make(map[int]*ObjectDefinitionSegment),
	}
}

type Reader struct {
	data        []byte
	cursor      int
	readLimit   int
	readLimited bool
}

func NewReader(data []byte) *Reader {
	return &Reader{
		data:        data,
		cursor:      0,
		readLimit:   len(data),
		readLimited: false,
	}
}

func (r *Reader) RemainingBytes() int {
	return len(r.data) - r.cursor
}

func (r *Reader) SetReadLimit(limit int) {
	r.readLimited = true
	r.readLimit = limit
}

func (r *Reader) RemoveReadLimit() {
	r.readLimited = false
	r.readLimit = r.RemainingBytes()
}

func (r *Reader) ReadLimit() int {
	return r.readLimit
}

func (r *Reader) SkipPastReadLimit() {
	if !r.readLimited {
		return
	}
	r.cursor += r.readLimit
	r.RemoveReadLimit()
}

func (r *Reader) ReachedEnd() bool {
	return r.cursor >= len(r.data)
}

func (r *Reader) Read(count int) ([]byte, error) {
	if count == 0 {
		return make([]byte, 0), nil
	}

	if r.readLimited {
		if count > r.readLimit {
			buf, err := r.Read(r.readLimit)
			if err != nil {
				return buf, err
			}

			return append(buf, make([]byte, count-len(buf))...), nil
		}

		r.readLimit -= count
	}

	r.cursor += count

	if r.cursor > len(r.data) {
		return make([]byte, 0), fmt.Errorf("Read cursor out of bound")
	}

	buf := make([]byte, count)
	copy(buf, r.data[r.cursor-count:r.cursor])

	return buf, nil
}

func (r *Reader) ReadByte() (byte, error) {
	bytes, err := r.Read(1)
	if err != nil {
		return 0x00, err
	}

	return bytes[0], nil
}
