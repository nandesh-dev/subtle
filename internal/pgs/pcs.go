package pgs

import (
	"encoding/binary"
	"fmt"
)

func ReadPresentationCompositionSegment(reader *Reader, header *Header) (*PresentationCompositionSegment, error) {
	reader.SetLimit(header.SegmentSize)
	defer reader.SkipPastLimit()

	buf, err := reader.Read(11)

	if err != nil {
		return nil, fmt.Errorf("Error reading data: %v", err)
	}

	width := int(
		binary.BigEndian.Uint16(buf[0:2]),
	)

	height := int(
		binary.BigEndian.Uint16(buf[2:4]),
	)

	compositionState, err := mapCompositionState(buf[7])
	if err != nil {
		return nil, err
	}

	paletteUpdateFlag, err := mapPaletteUpdateFlag(buf[8])
	if err != nil {
		return nil, err
	}

	paletteID := int(buf[9])

	compositionObjectCount := int(buf[10])
	compositionObjects := make([]CompositionObject, 0, compositionObjectCount)

	for len(compositionObjects) < compositionObjectCount {
		compositionObject, err := readCompositionObject(reader)
		if err != nil {
			return nil, err
		}

		compositionObjects = append(compositionObjects, *compositionObject)
	}

	return &PresentationCompositionSegment{
		Width:              width,
		Height:             height,
		CompositionState:   compositionState,
		PaletteUpdateFlag:  paletteUpdateFlag,
		PaletteID:          paletteID,
		CompositionObjects: compositionObjects,
	}, nil
}

func readCompositionObject(reader *Reader) (*CompositionObject, error) {
	buf, err := reader.Read(8)

	if err != nil {
		return nil, fmt.Errorf("Error reading composition object: %v", err)
	}

	objectID := int(
		binary.BigEndian.Uint16(buf[0:2]),
	)

	windowID := int(buf[2])

	objectCroppedFlag, err := mapObjectCroppedFlag(buf[3])
	if err != nil {
		return nil, err
	}

	objectHorizontalPosition := int(
		binary.BigEndian.Uint16(buf[4:6]),
	)

	objectVerticalPosition := int(
		binary.BigEndian.Uint16(buf[6:8]),
	)

	objectCroppingHorizontalPosition := 0
	objectCroppingVerticalPosition := 0
	objectCroppingWidth := 0
	objectCroppingHeight := 0

	if objectCroppedFlag {
		buf, err := reader.Read(8)

		if err != nil {
			return nil, fmt.Errorf("Error reading PCS composition object %v", err)
		}

		objectCroppingHorizontalPosition = int(
			binary.BigEndian.Uint16(buf[0:2]),
		)

		objectCroppingVerticalPosition = int(
			binary.BigEndian.Uint16(buf[2:4]),
		)

		objectCroppingWidth = int(
			binary.BigEndian.Uint16(buf[4:6]),
		)

		objectCroppingHeight = int(
			binary.BigEndian.Uint16(buf[6:8]),
		)
	}

	return &CompositionObject{
		ObjectID:                         objectID,
		WindowID:                         windowID,
		ObjectCroppedFlag:                objectCroppedFlag,
		ObjectHorizontalPosition:         objectHorizontalPosition,
		ObjectVerticalPosition:           objectVerticalPosition,
		ObjectCroppingHorizontalPosition: objectCroppingHorizontalPosition,
		ObjectCroppingVerticalPosition:   objectCroppingVerticalPosition,
		ObjectCroppingWidth:              objectCroppingWidth,
		ObjectCroppingHeight:             objectCroppingHeight,
	}, nil
}

func mapCompositionState(bt byte) (CompositionState, error) {
	switch bt {
	case 0x00:
		return Normal, nil
	case 0x40:
		return AcquisitionStart, nil
	case 0x80:
		return EpochStart, nil
	}

	return Normal, fmt.Errorf("Invalid composition state: %v", bt)
}

func mapPaletteUpdateFlag(bt byte) (bool, error) {
	switch bt {
	case 0x00:
		return false, nil
	case 0x80:
		return true, nil
	}

	return false, fmt.Errorf("Invalid 'palette update' flag: %v", bt)
}

func mapObjectCroppedFlag(bt byte) (bool, error) {
	switch bt {
	case 0x00:
		return false, nil
	case 0x40:
		return true, nil
	}

	return false, fmt.Errorf("Invalid 'object cropped' flag: %v", bt)
}
