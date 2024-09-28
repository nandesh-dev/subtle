package pgs

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/tesseract"
	"github.com/nandesh-dev/subtle/internal/warnings"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func DecodePGSSubtitle(rawStream subtitle.RawStream) (*subtitle.Subtitle, error, *warnings.WarningList) {
	warningList := warnings.NewWarningList()

	var subtitleBuf, errorBuf bytes.Buffer

	ffmpeg.LogCompiledCommand = false
	err := ffmpeg.Input(rawStream.VideoFilePath).
		Output("pipe:", ffmpeg.KwArgs{"map": fmt.Sprintf("0:%v", rawStream.Index), "c:s": "copy", "f": "sup"}).
		WithOutput(&subtitleBuf).
		WithErrorOutput(&errorBuf).
		Run()

	if err != nil {
		return nil, fmt.Errorf("Error extracting subtitles: %v %v", err, errorBuf), warningList
	}

	reader := NewReader(subtitleBuf.Bytes())

	displaySets := make([]displaySet, 0)
	dS := NewDisplaySet()

	for reader.RemainingBytes() > 11 {
		header, err := ReadHeader(reader)

		if err != nil {
			return nil, fmt.Errorf("Error reading header: %v", err), warningList
		}

		dS.Header = *header

		switch header.SegmentType {
		case PCS:
			segment, err := ReadPresentationCompositionSegment(reader, header)
			if err != nil {
				warningList.AddWarning(fmt.Errorf("Error reading PCS Segment: %v", err))

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
				warningList.AddWarning(fmt.Errorf("Error reading ODS segment: %v", err))
			} else {
				dS.ObjectDefinitionSegments[segment.ObjectID] = *segment
			}

		case PDS:
			segment, err := ReadPaletteDefinitionSegment(reader, header)
			if err != nil {
				warningList.AddWarning(fmt.Errorf("Error reading PDS segment: %v", err))
			} else {
				dS.PaletteDefinitionSegments[segment.PaletteID] = *segment
			}

		case WDS:
			segment, err := ReadWindowDefinitionSegment(reader, header)
			if err != nil {
				warningList.AddWarning(fmt.Errorf("Error reading WDS segment: %v", err))
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

	stream := subtitle.Stream{
		Langauge: rawStream.Language,
		Segments: make([]subtitle.Segment, 0),
	}

	tsrt := tesseract.NewClient()
	defer tsrt.Close()

	for _, displaySet := range displaySets {
		images, err := displaySet.parse()
		if err != nil {
			slog.Warn("Display set parsing error: %v %v", err, rawStream.VideoFilePath)
			continue
		}

		texts := make([]string, 0)

		for _, img := range images {
			text, err := tsrt.ExtractTextFromImage(img, stream.Langauge)
			if err != nil {
				return nil, err, warningList
			}

			if text != "" {
				texts = append(texts, text)
			}
		}

		segment := subtitle.Segment{
			Start: &displaySet.Header.PTS,
			End:   nil,
			Text:  strings.Join(texts, "\n"),
			Style: struct{}{},
		}

		stream.Segments = append(stream.Segments, segment)
	}

	return &subtitle.Subtitle{
		Streams: []subtitle.Stream{
			stream,
		},
	}, nil, warningList
}
