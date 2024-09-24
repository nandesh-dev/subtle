package decoder

import (
	"encoding/binary"
	"fmt"

	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/reader"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/segments"
)

func ReadObjectDefinitionSegment(reader *reader.Reader, header *segments.Header) (*segments.ObjectDefinitionSegment, error) {
	reader.SetLimit(header.SegmentSize)
	defer reader.SkipPastLimit()

	buf, err := reader.Read(7)

	if err != nil {
		return nil, fmt.Errorf("Error reading ODS %v", err)
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
		return nil, fmt.Errorf("Error reading ODS object data %v", err)
	}

	width := 0
	height := 0
	objectData := make([]byte, 0)

	if lastInSequenceFlag == segments.FirstInSequence || lastInSequenceFlag == segments.FirstAndLastInSequence {
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

	return &segments.ObjectDefinitionSegment{
		ObjectID:            objectID,
		ObjectVersionNumber: objectVersionNumber,
		LastInSquenceFlag:   lastInSequenceFlag,
		Width:               width,
		Height:              height,
		ObjectData:          objectData,
	}, nil
}

func mapLastInSequenceFlag(bt byte) (segments.LastInSquenceFlag, error) {
	switch bt {
	case 0x40:
		return segments.LastInSequence, nil
	case 0x80:
		return segments.FirstInSequence, nil
	case 0xC0:
		return segments.FirstAndLastInSequence, nil
	}

	return segments.FirstInSequence, fmt.Errorf("Invalid ODS law is sequence flag %v", bt)
}
