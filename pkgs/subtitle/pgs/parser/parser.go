package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"time"

	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/otiai10/gosseract/v2"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (_ *Parser) Parse(data []byte) (*subtitle.Subtitle, error) {
	tesseractClient := gosseract.NewClient()
	defer tesseractClient.Close()

	sub := subtitle.Subtitle{}

	reader := NewReader(data)

	previousDisplaySet := *NewDisplaySet()
	currentDisplaySet := *NewDisplaySet()

	for reader.RemainingBytes() > 11 {
		header, err := readHeader(reader)
		if err != nil {
			return nil, fmt.Errorf("Cannot reader header: %w", err)
		}

		currentDisplaySet.header = *header

		switch header.segmentType {
		case PCSSegment:
			segment, err := readPresentationCompositionSegment(reader, header)
			if err != nil {
				continue
			}

			if segment.state == AcquisitionStartPresentationCompositionState || (segment.state == NormalPresentationCompositionState && len(segment.objects) != 0) {
				for id, objectDefinitionSegment := range previousDisplaySet.objectDefinitionSegments {
					currentDisplaySet.objectDefinitionSegments[id] = objectDefinitionSegment
				}

				for id, windowDefinition := range previousDisplaySet.windowDefinitions {
					currentDisplaySet.windowDefinitions[id] = windowDefinition
				}

				for id, paletteDefinitionSegment := range previousDisplaySet.paletteDefinitionSegments {
					currentDisplaySet.paletteDefinitionSegments[id] = paletteDefinitionSegment
				}
			}

			currentDisplaySet.presentationCompositionSegment = *segment
		case ODSSegment:
			segment, err := readObjectDefinitionSegment(reader, header)
			if err != nil {
				continue
			}

			currentDisplaySet.objectDefinitionSegments[segment.objectId] = *segment

		case PDSSegment:
			segment, err := readPaletteDefinitionSegment(reader, header)
			if err != nil {
				continue
			}

			currentDisplaySet.paletteDefinitionSegments[segment.paletteId] = *segment

		case WDSSegment:
			segment, err := readWindowDefinitionSegment(reader, header)
			if err != nil {
				continue
			}

			for _, window := range segment.windows {
				currentDisplaySet.windowDefinitions[window.id] = window
			}

		case ENDSegment:
			images, err := extractImagesFromDisplaySet(currentDisplaySet)
			if err != nil {
				return nil, fmt.Errorf("Cannot extract images from display set: %w", err)
			}

			content := make([]subtitle.CueContentSegment, 0)

			var imageBuffer bytes.Buffer
			for _, image := range images {
				imageBuffer.Reset()

				if err := png.Encode(&imageBuffer, image); err != nil {
					return nil, fmt.Errorf("Cannot encode image to png: %w", err)
				}

				if err := tesseractClient.SetImageFromBytes(imageBuffer.Bytes()); err != nil {
					return nil, fmt.Errorf("Cannot send image to tesseract: %w", err)
				}

				text, err := tesseractClient.Text()
				if err != nil {
					return nil, fmt.Errorf("Cannot get text from tesseract: %w", err)
				}

				content = append(content, subtitle.CueContentSegment{
					Text: text,
				})
			}

			cue := subtitle.Cue{
				Timestamp: subtitle.CueTimestamp{
					Start: currentDisplaySet.header.pts,
				},
				OriginalImages: images,
				Content:        content,
			}

			sub.Cues = append(sub.Cues, cue)

			if len(sub.Cues) >= 2 {
				sub.Cues[len(sub.Cues)-2].Timestamp.End = currentDisplaySet.header.pts
			}

			previousDisplaySet = currentDisplaySet
			currentDisplaySet = *NewDisplaySet()
		}
	}

	if len(sub.Cues) >= 1 {
		sub.Cues[len(sub.Cues)-1].Timestamp.End = sub.Cues[len(sub.Cues)-1].Timestamp.Start + 15*time.Second
	}

	return &sub, nil
}

func readHeader(reader *Reader) (*Header, error) {
	reader.SetReadLimit(13)
	defer reader.RemoveReadLimit()

	buf, err := reader.Read(11)
	if err != nil {
		return nil, fmt.Errorf("Error reading bytes: %w", err)
	}

	if !bytes.Equal(buf[0:2], []byte{0x50, 0x47}) {
		return nil, fmt.Errorf("Incorrect magic number: %v", buf[0:2])
	}

	pts := time.Duration(
		int(binary.BigEndian.Uint32(buf[2:6])) * 1_000_000 / 90,
	)

	var segmentType SegmentType

	switch buf[10] {
	case 0x14:
		segmentType = PDSSegment
	case 0x15:
		segmentType = ODSSegment
	case 0x16:
		segmentType = PCSSegment
	case 0x17:
		segmentType = WDSSegment
	case 0x80:
		segmentType = ENDSegment
	default:
		return nil, fmt.Errorf("Invalid segment type byte: %v", buf[10])
	}

	segmentSize := 0

	if reader.RemainingBytes() >= 2 {
		buf, err = reader.Read(2)
		if err != nil {
			return nil, fmt.Errorf("Cannot read segment size: %w", err)
		}

		segmentSize = int(
			binary.BigEndian.Uint16(buf),
		)
	}

	return &Header{
		pts:         pts,
		segmentType: segmentType,
		segmentSize: segmentSize,
	}, nil
}

func readPresentationCompositionSegment(reader *Reader, header *Header) (*PresentationCompositionSegment, error) {
	reader.SetReadLimit(header.segmentSize)
	defer reader.SkipPastReadLimit()

	segment := PresentationCompositionSegment{}

	buf, err := reader.Read(11)
	if err != nil {
		return nil, fmt.Errorf("Error reading data: %w", err)
	}

	segment.width = int(
		binary.BigEndian.Uint16(buf[0:2]),
	)

	segment.height = int(
		binary.BigEndian.Uint16(buf[2:4]),
	)

	switch buf[7] {
	case 0x00:
		segment.state = NormalPresentationCompositionState
	case 0x40:
		segment.state = AcquisitionStartPresentationCompositionState
	case 0x80:
		segment.state = EpochStartPresentationCompositionState
	default:
		return nil, fmt.Errorf("Invalid composition state: %v", buf[7])
	}

	switch buf[8] {
	case 0x00:
		segment.paletteUpdateFlag = false
	case 0x80:
		segment.paletteUpdateFlag = true
	default:
		return nil, fmt.Errorf("Invalid palette update flag: %v", buf[8])
	}

	segment.paletteId = int(buf[9])

	compositionObjectCount := int(buf[10])

	for i := 0; i < compositionObjectCount; i++ {
		compositionObject := PresentationCompositionObject{}

		buf, err := reader.Read(8)
		if err != nil {
			return nil, fmt.Errorf("Cannot read composition object bytes: %w", err)
		}

		compositionObject.objectId = int(
			binary.BigEndian.Uint16(buf[0:2]),
		)

		compositionObject.windowId = int(buf[2])

		switch buf[3] {
		case 0x00:
			compositionObject.objectCroppedFlag = false
		case 0x40:
			compositionObject.objectCroppedFlag = true
		default:
			return nil, fmt.Errorf("Invalid object cropped flag: %v", buf[3])
		}

		compositionObject.objectHorizontalPosition = int(
			binary.BigEndian.Uint16(buf[4:6]),
		)

		compositionObject.objectVerticalPosition = int(
			binary.BigEndian.Uint16(buf[6:8]),
		)

		if compositionObject.objectCroppedFlag {
			buf, err := reader.Read(8)
			if err != nil {
				return nil, fmt.Errorf("Cannot read object cropping dimensions: %w", err)
			}

			compositionObject.objectCroppingHorizontalPosition = int(
				binary.BigEndian.Uint16(buf[0:2]),
			)

			compositionObject.objectCroppingVerticalPosition = int(
				binary.BigEndian.Uint16(buf[2:4]),
			)

			compositionObject.objectCroppingWidth = int(
				binary.BigEndian.Uint16(buf[4:6]),
			)

			compositionObject.objectCroppingHeight = int(
				binary.BigEndian.Uint16(buf[6:8]),
			)
		}

		segment.objects = append(segment.objects, compositionObject)
	}

	return &segment, nil
}

func readObjectDefinitionSegment(reader *Reader, header *Header) (*ObjectDefinitionSegment, error) {
	reader.SetReadLimit(header.segmentSize)
	defer reader.SkipPastReadLimit()

	objectDefinitionSegment := ObjectDefinitionSegment{}

	buf, err := reader.Read(7)
	if err != nil {
		return nil, fmt.Errorf("Cannot read data: %w", err)
	}

	objectDefinitionSegment.objectId = int(
		binary.BigEndian.Uint16(buf[0:2]),
	)

	objectDefinitionSegment.objectVersionNumber = int(buf[2])

	switch buf[3] {
	case 0x40:
		objectDefinitionSegment.sequenceFlag = LastInObjectDefinitionSequence
	case 0x80:
		objectDefinitionSegment.sequenceFlag = FirstInObjectDefinitionSequence
	case 0xC0:
		objectDefinitionSegment.sequenceFlag = FirstAndLastInObjectDefinitionSequence
	default:
		return nil, fmt.Errorf("Invalid sequence flag: %v", buf[3])
	}

	objectDataLengthBytes := append([]byte{0x00}, buf[4:7]...)
	objectDataLength := int(
		binary.BigEndian.Uint32(objectDataLengthBytes),
	)

	buf, err = reader.Read(objectDataLength)
	if err != nil {
		return nil, fmt.Errorf("Cannot read object data: %w", err)
	}

	objectDefinitionSegment.objectData = buf[:]

	if objectDefinitionSegment.sequenceFlag == FirstInObjectDefinitionSequence || objectDefinitionSegment.sequenceFlag == FirstAndLastInObjectDefinitionSequence {
		objectDefinitionSegment.width = int(
			binary.BigEndian.Uint16(buf[0:2]),
		)

		objectDefinitionSegment.height = int(
			binary.BigEndian.Uint16(buf[2:4]),
		)

		objectDefinitionSegment.objectData = buf[4:]
	}

	return &objectDefinitionSegment, nil
}

func readPaletteDefinitionSegment(reader *Reader, header *Header) (*PaletteDefinitionSegment, error) {
	reader.SetReadLimit(header.segmentSize)
	defer reader.SkipPastReadLimit()

	segment := PaletteDefinitionSegment{
		paletteEntries: make(map[int]color.Color, 0),
	}

	rawPaletteID, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("Cannot read palette id: %w", err)
	}
	segment.paletteId = int(rawPaletteID)

	rawPaletteVersionNumber, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("Cannot read palette version number: %w", err)
	}
	segment.paletteVersionNumber = int(rawPaletteVersionNumber)

	paletteEntriesCount := (header.segmentSize - 2) / 5

	clamp := func(number float64, min int, max int) int {
		return int(math.Max(float64(min), math.Min(float64(max), number)))
	}

	for len(segment.paletteEntries) < paletteEntriesCount {
		buf, err := reader.Read(5)

		if err != nil {
			return nil, fmt.Errorf("Cannot read palette entry: %w", err)
		}

		paletteEntryId := int(buf[0])

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

		segment.paletteEntries[paletteEntryId] = paletteColor
	}

	return &segment, nil
}

func readWindowDefinitionSegment(reader *Reader, header *Header) (*WindowDefinitionSegment, error) {
	reader.SetReadLimit(header.segmentSize)
	defer reader.SkipPastReadLimit()

	segment := WindowDefinitionSegment{}

	rawNumberOfWindows, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("Error reading WDS number of windows %v", err)
	}
	numberOfWindows := int(rawNumberOfWindows)

	for len(segment.windows) < numberOfWindows {
		buf, err := reader.Read(9)
		if err != nil {
			return nil, fmt.Errorf("Error reading WDS %v", err)
		}

		window := Window{}

		window.id = int(buf[0])
		window.horizontalPosition = int(
			binary.BigEndian.Uint16(buf[1:3]),
		)
		window.verticalPosition = int(
			binary.BigEndian.Uint16(buf[3:5]),
		)
		window.width = int(
			binary.BigEndian.Uint16(buf[5:7]),
		)
		window.height = int(
			binary.BigEndian.Uint16(buf[7:9]),
		)

		segment.windows = append(segment.windows, window)
	}

	return &segment, nil
}

func extractImagesFromDisplaySet(displaySet DisplaySet) ([]image.Image, error) {
	if len(displaySet.presentationCompositionSegment.objects) == 0 {
		return nil, nil
	}

	paletteId := displaySet.presentationCompositionSegment.paletteId
	paletteDefinitionSegment, exist := displaySet.paletteDefinitionSegments[paletteId]
	if !exist {
		return nil, fmt.Errorf("Palette doesn't exist with id: %v", paletteId)
	}

	images := make([]image.Image, 0)

	for _, compositionObject := range displaySet.presentationCompositionSegment.objects {
		objectId := compositionObject.objectId
		objectDefinitionSegment, exist := displaySet.objectDefinitionSegments[objectId]
		if !exist {
			return nil, fmt.Errorf("Object doesn't exist with id: %v", objectId)
		}

		reader := NewReader(objectDefinitionSegment.objectData)

		img := image.NewRGBA(image.Rect(0, 0, objectDefinitionSegment.width, objectDefinitionSegment.height))
		xCoord, yCoord := 0, 0

		fillImage := func(id int, count int) {
			c, exist := paletteDefinitionSegment.paletteEntries[id]
			if !exist {
				c = color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 255,
				}
			}

			for i := 0; i < count; i++ {
				img.Set(xCoord, yCoord, c)
				xCoord++
			}
		}

		for !reader.ReachedEnd() {
			firstByte, err := reader.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("Cannot read first object definition byte: %w", err)
			}

			if firstByte > 0 {
				fillImage(int(firstByte), 1)
				continue
			}

			secondByte, err := reader.ReadByte()

			if err != nil {
				return nil, fmt.Errorf("Cannot read second object definition byte: %w", err)
			}

			if secondByte == 0 {
				yCoord++
				xCoord = 0
				continue
			}

			if secondByte < 64 {
				fillImage(0, 2)
				continue
			}

			thirdByte, err := reader.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("Cannot read third object definition byte: %w", err)
			}

			if secondByte < 128 {
				count := int((secondByte-64))*256 + int(thirdByte)
				fillImage(0, count)
				continue
			}

			if secondByte < 192 {
				count := int(secondByte) - 128
				fillImage(int(thirdByte), count)
				continue
			}

			forthByte, err := reader.ReadByte()

			if err != nil {
				return nil, fmt.Errorf("Cannot read forth object definition byte: %w", err)
			}

			count := int(secondByte-192)*256 + int(thirdByte)
			fillImage(int(forthByte), count)
		}

		images = append(images, img)
	}

	return images, nil
}
