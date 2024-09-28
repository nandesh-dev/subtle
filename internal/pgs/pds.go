package pgs

import (
	"fmt"
	"image/color"
	"math"
)

func ReadPaletteDefinitionSegment(reader *Reader, header *Header) (*PaletteDefinitionSegment, error) {
	reader.SetLimit(header.SegmentSize)
	defer reader.SkipPastLimit()

	rawPaletteID, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("Error reading palette ID: %v", err)
	}

	paletteID := int(rawPaletteID)

	rawPaletteVersionNumber, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("Error reading pallete version number: %v", err)
	}

	paletteVersionNumber := int(rawPaletteVersionNumber)

	paletteEntriesCount := (header.SegmentSize - 2) / 5

	paletteEntries := map[int]color.Color{}

	for len(paletteEntries) < paletteEntriesCount {
		buf, err := reader.Read(5)

		if err != nil {
			return nil, fmt.Errorf("Error reading palette entry: %v", err)
		}

		paletteEntryID := int(buf[0])

		y := float64(buf[1])
		cr := float64(buf[2])
		cb := float64(buf[3])
		transparency := int(buf[4])

		r := clamp(math.Floor(y+1.4075*(cr-128)), 0, 255)
		g := clamp(math.Floor(y-0.3455*(cb-128)-0.7169*(cr-128)), 0, 255)
		b := clamp(math.Floor(y+1.779*(cb-128)), 0, 255)

		paletteColor := color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(transparency),
		}

		paletteEntries[paletteEntryID] = paletteColor
	}

	return &PaletteDefinitionSegment{
		PaletteID:            paletteID,
		PaletteVersionNumber: paletteVersionNumber,
		PaletteEntries:       paletteEntries,
	}, nil
}

func clamp(number float64, min int, max int) int {
	return int(math.Max(float64(min), math.Min(float64(max), number)))
}
