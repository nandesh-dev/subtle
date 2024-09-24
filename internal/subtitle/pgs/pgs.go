package pgs

import (
	"fmt"
	"log/slog"

	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/decoder"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/reader"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/segments"
)

func DecodePGSSubtitle(data []byte) (*subtitle.Subtitle, error) {
	reader := reader.NewReader(data)

	displaySets := make([]displaySet, 0)
	currentDisplaySet := displaySet{
		PaletteDefinitionSegments: make(map[int]*segments.PaletteDefinitionSegment),
		WindowDefinitionSegments:  make(map[int]*segments.Window),
		ObjectDefinitionSegments:  make(map[int]*segments.ObjectDefinitionSegment),
	}
	previousDisplaySet := displaySet{
		PaletteDefinitionSegments: make(map[int]*segments.PaletteDefinitionSegment),
		WindowDefinitionSegments:  make(map[int]*segments.Window),
		ObjectDefinitionSegments:  make(map[int]*segments.ObjectDefinitionSegment),
	}

	for reader.RemainingBytes() > 11 {
		header, err := decoder.ReadHeader(reader)

		currentDisplaySet.Header = header

		if err != nil {
			return nil, err
		}

		switch header.SegmentType {
		case segments.PCS:
			segment, err := decoder.ReadPresentationCompositionSegment(reader, header)
			if err != nil {
				slog.Warn("PGS Decoder PCS Decoding Error", "error", err)

			} else {
				if len(displaySets) >= 0 && (segment.CompositionState == segments.AcquisitionStart || (segment.CompositionState == segments.Normal && len(segment.CompositionObjects) != 0)) {
					for id, object := range previousDisplaySet.ObjectDefinitionSegments {
						currentDisplaySet.ObjectDefinitionSegments[id] = object
					}

					for id, window := range previousDisplaySet.WindowDefinitionSegments {
						currentDisplaySet.WindowDefinitionSegments[id] = window
					}

					for id, palette := range previousDisplaySet.WindowDefinitionSegments {
						currentDisplaySet.WindowDefinitionSegments[id] = palette
					}
				}

				currentDisplaySet.PresentationCompositionSegment = segment
			}

		case segments.ODS:
			segment, err := decoder.ReadObjectDefinitionSegment(reader, header)
			if err != nil {
				slog.Warn("PGS Decoder ODS Decoding Error", "error", err)
			} else {
				currentDisplaySet.ObjectDefinitionSegments[segment.ObjectID] = segment
			}

		case segments.PDS:
			segment, err := decoder.ReadPaletteDefinitionSegment(reader, header)
			if err != nil {
				slog.Warn("PGS Decoder PDS Decoding Error", "error", err)
			} else {
				currentDisplaySet.PaletteDefinitionSegments[segment.PaletteID] = segment
			}

		case segments.WDS:
			segment, err := decoder.ReadWindowDefinitionSegment(reader, header)
			if err != nil {
				slog.Warn("PGS Decoder WDS Decoding Error", "error", err)
			} else {
				for _, window := range segment.Windows {
					currentDisplaySet.WindowDefinitionSegments[window.WindowID] = &window
				}
			}

		case segments.END:
			displaySets = append(displaySets, currentDisplaySet)
			previousDisplaySet = currentDisplaySet
			currentDisplaySet = displaySet{
				PaletteDefinitionSegments: make(map[int]*segments.PaletteDefinitionSegment),
				WindowDefinitionSegments:  make(map[int]*segments.Window),
				ObjectDefinitionSegments:  make(map[int]*segments.ObjectDefinitionSegment),
			}
		}
	}

	stream, err := parseDisplaySets(displaySets)

	if err != nil {
		return nil, fmt.Errorf("Error paring display set: %v", err)
	}

	return stream, nil
}
