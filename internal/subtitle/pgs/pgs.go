package pgs

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"github.com/nandesh-dev/subtle/internal/ocr/tesseract"
	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/decoder"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/reader"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/segments"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func DecodePGSSubtitle(rawStream subtitle.RawSubtitleStream) (*subtitle.Subtitle, error) {
	var subtitleBuf, errorBuf bytes.Buffer

	ffmpeg.LogCompiledCommand = false
	err := ffmpeg.Input(rawStream.VideoFilePath).
		Output("pipe:", ffmpeg.KwArgs{"map": fmt.Sprintf("0:%v", rawStream.Index), "c:s": "copy", "f": "sup"}).
		WithOutput(&subtitleBuf).
		WithErrorOutput(&errorBuf).
		Run()

	if err != nil {
		return nil, fmt.Errorf("Error extracting subtitles: %v %v", err, errorBuf)
	}

	reader := reader.NewReader(subtitleBuf.Bytes())

	displaySets := make([]displaySet, 0)
	dS := NewDisplaySet()

	for reader.RemainingBytes() > 11 {
		header, err := decoder.ReadHeader(reader)

		dS.Header = header

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
					for id, object := range displaySets[len(displaySets)-1].ObjectDefinitionSegments {
						dS.ObjectDefinitionSegments[id] = object
					}

					for id, window := range displaySets[len(displaySets)-1].WindowDefinitionSegments {
						dS.WindowDefinitionSegments[id] = window
					}

					for id, palette := range displaySets[len(displaySets)-1].WindowDefinitionSegments {
						dS.WindowDefinitionSegments[id] = palette
					}
				}

				dS.PresentationCompositionSegment = segment
			}

		case segments.ODS:
			segment, err := decoder.ReadObjectDefinitionSegment(reader, header)
			if err != nil {
				slog.Warn("PGS Decoder ODS Decoding Error", "error", err)
			} else {
				dS.ObjectDefinitionSegments[segment.ObjectID] = segment
			}

		case segments.PDS:
			segment, err := decoder.ReadPaletteDefinitionSegment(reader, header)
			if err != nil {
				slog.Warn("PGS Decoder PDS Decoding Error", "error", err)
			} else {
				dS.PaletteDefinitionSegments[segment.PaletteID] = segment
			}

		case segments.WDS:
			segment, err := decoder.ReadWindowDefinitionSegment(reader, header)
			if err != nil {
				slog.Warn("PGS Decoder WDS Decoding Error", "error", err)
			} else {
				for _, window := range segment.Windows {
					dS.WindowDefinitionSegments[window.WindowID] = &window
				}
			}

		case segments.END:
			displaySets = append(displaySets, dS)
			dS = NewDisplaySet()
		}
	}

	stream := subtitle.SubtitleStream{
		Langauge: rawStream.Language,
		Segments: make([]subtitle.SubtitleSegment, 0),
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
				return nil, err
			}

			if text != "" {
				texts = append(texts, text)
			}
		}

		segment := subtitle.SubtitleSegment{
			Start: &displaySet.Header.PTS,
			End:   nil,
			Text:  strings.Join(texts, "\n"),
			Style: struct{}{},
		}

		stream.Segments = append(stream.Segments, segment)
	}

	return &subtitle.Subtitle{
		Streams: []subtitle.SubtitleStream{
			stream,
		},
	}, nil
}
