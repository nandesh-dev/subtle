package pgs

import (
	"encoding/binary"
	"fmt"
)

func ReadObjectDefinitionSegment(reader *Reader, header *Header) (*ObjectDefinitionSegment, error) {
	reader.SetLimit(header.SegmentSize)
	defer reader.SkipPastLimit()

	buf, err := reader.Read(7)

	if err != nil {
		return nil, fmt.Errorf("Error reading data: %v", err)
	}

	objectID := int(
		binary.BigEndian.Uint16(buf[0:2]),
	)

	objectVersionNumber := int(buf[2])

	lastInSequenceFlag, err := mapLastInSequenceFlag(buf[3])
	if err != nil {
		return nil, err
	}

	objectDataLengthBytes := append([]byte{0x00}, buf[4:7]...)
	objectDataLength := int(
		binary.BigEndian.Uint32(objectDataLengthBytes),
	)

	buf, err = reader.Read(objectDataLength)
	if err != nil {
		return nil, fmt.Errorf("Error reading object data: %v", err)
	}

	width := 0
	height := 0
	objectData := make([]byte, 0)

	if lastInSequenceFlag == FirstInSequence || lastInSequenceFlag == FirstAndLastInSequence {
		width = int(
			binary.BigEndian.Uint16(buf[0:2]),
		)

		height = int(
			binary.BigEndian.Uint16(buf[2:4]),
		)

		objectData = buf[4:]
	} else {
		objectData = buf[:]
	}

	return &ObjectDefinitionSegment{
		ObjectID:            objectID,
		ObjectVersionNumber: objectVersionNumber,
		LastInSquenceFlag:   lastInSequenceFlag,
		Width:               width,
		Height:              height,
		ObjectData:          objectData,
	}, nil
}

func mapLastInSequenceFlag(bt byte) (LastInSquenceFlag, error) {
	switch bt {
	case 0x40:
		return LastInSequence, nil
	case 0x80:
		return FirstInSequence, nil
	case 0xC0:
		return FirstAndLastInSequence, nil
	}

	return FirstInSequence, fmt.Errorf("Invalid 'last in sequence' flag: %v", bt)
}
