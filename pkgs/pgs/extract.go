package pgs

import (
	"bytes"
	"fmt"
	"slices"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/warning"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ExtractFromRawStream(rawStream filemanager.RawStream) (*subtitle.ImageSubtitle, warning.WarningList, error) {
	warnings := warning.NewWarningList()

	var subtitleBuf, errorBuf bytes.Buffer

	ffmpeg.LogCompiledCommand = false
	err := ffmpeg.Input(rawStream.Filepath()).
		Output("pipe:", ffmpeg.KwArgs{"map": fmt.Sprintf("0:%v", rawStream.Index()), "c:s": "copy", "f": "sup"}).
		WithOutput(&subtitleBuf).
		WithErrorOutput(&errorBuf).
		Run()

	if err != nil {
		return nil, *warnings, fmt.Errorf("Error extracting subtitles: %v %v", err, errorBuf)
	}

	reader := NewReader(subtitleBuf.Bytes())

	displaySets := make([]displaySet, 0)
	dS := NewDisplaySet()

	for reader.RemainingBytes() > 11 {
		header, err := ReadHeader(reader)

		if err != nil {
			return nil, *warnings, fmt.Errorf("Error reading header: %v", err)
		}

		dS.Header = *header

		switch header.SegmentType {
		case PCS:
			segment, err := ReadPresentationCompositionSegment(reader, header)
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error reading PCS Segment: %v", err))

			} else {
				if len(displaySets) >= 0 && (segment.CompositionState == AcquisitionStart || (segment.CompositionState == Normal && len(segment.CompositionObjects) != 0)) {
					for id, object := range displaySets[len(displaySets)-1].ObjectDefinitionSegments {
						dS.ObjectDefinitionSegments[id] = object
					}

					for id, window := range displaySets[len(displaySets)-1].WindowDefinitions {
						dS.WindowDefinitions[id] = window
					}

					for id, palette := range displaySets[len(displaySets)-1].PaletteDefinitionSegments {
						dS.PaletteDefinitionSegments[id] = palette
					}
				}

				dS.PresentationCompositionSegment = *segment
			}

		case ODS:
			segment, err := ReadObjectDefinitionSegment(reader, header)
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error reading ODS segment: %v", err))
			} else {
				dS.ObjectDefinitionSegments[segment.ObjectID] = *segment
			}

		case PDS:
			segment, err := ReadPaletteDefinitionSegment(reader, header)
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error reading PDS segment: %v", err))
			} else {
				dS.PaletteDefinitionSegments[segment.PaletteID] = *segment
			}

		case WDS:
			segment, err := ReadWindowDefinitionSegment(reader, header)
			if err != nil {
				warnings.AddWarning(fmt.Errorf("Error reading WDS segment: %v", err))
			} else {
				for _, window := range segment.Windows {
					dS.WindowDefinitions[window.WindowID] = window
				}
			}

		case END:
			displaySets = append(displaySets, dS)
			dS = NewDisplaySet()
		}
	}

	sub := subtitle.NewImageSubtitle()

	previousStartTimestamp := time.Second * 0

	slices.Reverse(displaySets)

	for _, displaySet := range displaySets {
		images, err := displaySet.parse()
		if err != nil {
			warnings.AddWarning(fmt.Errorf("Display set parsing error: %v", err))
			continue
		}

		for _, image := range images {
			endTimestamp := previousStartTimestamp
			if previousStartTimestamp.Nanoseconds() < displaySet.Header.PTS.Nanoseconds() {
				endTimestamp = displaySet.Header.PTS + time.Second*15
			}

			segment := subtitle.NewImageSegment(displaySet.Header.PTS, endTimestamp, image)
			sub.AddSegment(*segment)
		}

		previousStartTimestamp = displaySet.Header.PTS
	}

	return sub, *warnings, nil
}
