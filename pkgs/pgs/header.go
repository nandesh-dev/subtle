package pgs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type SegmentType int

const (
	PDS SegmentType = iota
	ODS
	PCS
	WDS
	END
)

type Header struct {
	PTS         time.Duration
	SegmentType SegmentType
	SegmentSize int
}

func mapSegmentType(bt byte) (SegmentType, error) {
	switch bt {
	case 0x14:
		return PDS, nil
	case 0x15:
		return ODS, nil
	case 0x16:
		return PCS, nil
	case 0x17:
		return WDS, nil
	case 0x80:
		return END, nil
	}

	return END, fmt.Errorf("Invalid segment type: %v", bt)
}

func ReadHeader(reader *Reader) (*Header, error) {
	reader.SetLimit(13)

	defer reader.RemoveLimit()

	buf, err := reader.Read(11)

	if err != nil {
		return nil, fmt.Errorf("Error reading header %v", err)
	}

	if !bytes.Equal(buf[0:2], []byte{0x50, 0x47}) {
		return nil, fmt.Errorf("Incorrect magic number %v", buf[0:2])
	}

	pts := time.Duration(
		int(binary.BigEndian.Uint32(buf[2:6])) * 1000 / 90,
	)

	segmentType, err := mapSegmentType(buf[10])
	if err != nil {
		return nil, err
	}

	segmentSize := 0

	if reader.RemainingBytes() >= 2 {
		buf, err = reader.Read(2)
		if err != nil {
			return nil, fmt.Errorf("Error reading header segment size: %v", err)
		}

		segmentSize = int(
			binary.BigEndian.Uint16(buf),
		)
	}

	return &Header{
		PTS:         pts,
		SegmentType: segmentType,
		SegmentSize: segmentSize,
	}, nil
}
