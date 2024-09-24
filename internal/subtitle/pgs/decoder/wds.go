package decoder

import (
	"encoding/binary"
	"fmt"

	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/reader"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/segments"
)

func ReadWindowDefinitionSegment(reader *reader.Reader, header *segments.Header) (*segments.WindowDefinitionSegment, error) {
	reader.SetLimit(header.SegmentSize)
	defer reader.SkipPastLimit()

	rawNumberOfWindows, err := reader.ReadByte()

	if err != nil {
		return nil, fmt.Errorf("Error reading WDS number of windows %v", err)
	}

	numberOfWindows := int(rawNumberOfWindows)

	windows := make([]segments.Window, 0, numberOfWindows)

	for len(windows) < numberOfWindows {
		buf, err := reader.Read(9)

		if err != nil {
			return nil, fmt.Errorf("Error reading WDS %v", err)
		}

		windowID := int(buf[0])

		windowHorizontalPosition := int(
			binary.BigEndian.Uint16(buf[1:3]),
		)

		windowVerticalPosition := int(
			binary.BigEndian.Uint16(buf[3:5]),
		)

		windowWidth := int(
			binary.BigEndian.Uint16(buf[5:7]),
		)

		windowHeight := int(
			binary.BigEndian.Uint16(buf[7:9]),
		)

		window := segments.Window{
			WindowID:                 windowID,
			WindowHorizontalPosition: windowHorizontalPosition,
			WindowVerticalPosition:   windowVerticalPosition,
			WindowWidth:              windowWidth,
			WindowHeight:             windowHeight,
		}

		windows = append(windows, window)
	}

	return &segments.WindowDefinitionSegment{
		Windows: windows,
	}, nil
}
