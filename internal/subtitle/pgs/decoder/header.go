package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/reader"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/segments"
)

func ReadHeader(reader *reader.Reader) (*segments.Header, error) {
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

	segmentType := segments.INV

	switch buf[10] {
	case 0x14:
		segmentType = segments.PDS
	case 0x15:
		segmentType = segments.ODS
	case 0x16:
		segmentType = segments.PCS
	case 0x17:
		segmentType = segments.WDS
	case 0x80:
		segmentType = segments.END
	}

	if segmentType == segments.INV {
		return nil, fmt.Errorf("Invalid segment type %v", buf[10])
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

	return &segments.Header{
		PTS:         pts,
		SegmentType: segmentType,
		SegmentSize: segmentSize,
	}, nil
}
