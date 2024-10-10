package pgs

import (
	"encoding/binary"
	"fmt"
)

func ReadWindowDefinitionSegment(reader *Reader, header *Header) (*WindowDefinitionSegment, error) {
	reader.SetLimit(header.SegmentSize)
	defer reader.SkipPastLimit()

	rawNumberOfWindows, err := reader.ReadByte()

	if err != nil {
		return nil, fmt.Errorf("Error reading WDS number of windows %v", err)
	}

	numberOfWindows := int(rawNumberOfWindows)

	windows := make([]Window, 0, numberOfWindows)

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

		window := Window{
			WindowID:                 windowID,
			WindowHorizontalPosition: windowHorizontalPosition,
			WindowVerticalPosition:   windowVerticalPosition,
			WindowWidth:              windowWidth,
			WindowHeight:             windowHeight,
		}

		windows = append(windows, window)
	}

	return &WindowDefinitionSegment{
		Windows: windows,
	}, nil
}
