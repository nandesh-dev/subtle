package pgs

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/reader"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs/segments"
)

type displaySet struct {
	Header                         *segments.Header
	PresentationCompositionSegment *segments.PresentationCompositionSegment
	WindowDefinitionSegments       map[int]*segments.Window
	PaletteDefinitionSegments      map[int]*segments.PaletteDefinitionSegment
	ObjectDefinitionSegments       map[int]*segments.ObjectDefinitionSegment
}

func NewDisplaySet() displaySet {
	return displaySet{
		PaletteDefinitionSegments: make(map[int]*segments.PaletteDefinitionSegment),
		WindowDefinitionSegments:  make(map[int]*segments.Window),
		ObjectDefinitionSegments:  make(map[int]*segments.ObjectDefinitionSegment),
	}
}

func (d *displaySet) parse() ([]image.Image, error) {
	if len(d.PresentationCompositionSegment.CompositionObjects) == 0 {
		return nil, nil
	}

	paletteID := d.PresentationCompositionSegment.PaletteID
	palette, paletteExist := d.PaletteDefinitionSegments[paletteID]

	if !paletteExist {
		return nil, fmt.Errorf("Palette not found")
	}

	images := make([]image.Image, 0)

	for _, compositionObject := range d.PresentationCompositionSegment.CompositionObjects {
		objectID := compositionObject.ObjectID
		object, objectExist := d.ObjectDefinitionSegments[objectID]

		if !objectExist {
			continue
		}

		imageColorIDs, err := decodeRLEImageData(object.ObjectData)

		if err != nil {
			log.Fatal(err)
		}

		img := image.NewRGBA(image.Rect(0, 0, object.Width, object.Height))

		for y, line := range imageColorIDs {
			for x, colorID := range line {
				if colorID > len(palette.PaletteEntries) {
					img.Set(x, y, color.RGBA{
						R: 0,
						G: 0,
						B: 0,
						A: 255,
					})
				} else {
					color, colorExist := palette.PaletteEntries[colorID]

					if !colorExist {
						continue
					}

					img.Set(x, y, color)
				}
			}
		}

		images = append(images, img)
	}

	return images, nil
}

func decodeRLEImageData(data []byte) ([][]int, error) {

	reader := reader.NewReader(data)

	imageColorIDs := make([][]int, 0)
	currentLine := make([]int, 0)

	addIDMultipleTimes := func(id int, count int) {
		ids := make([]int, count)

		for i := range ids {
			ids[i] = id
		}

		currentLine = append(currentLine, ids...)
	}

	for !reader.ReachedEnd() {
		firstByte, err := reader.ReadByte()

		if err != nil {
			return make([][]int, 0), fmt.Errorf("Error reading first byte: %v", err)
		}

		if firstByte > 0 {
			currentLine = append(currentLine, int(firstByte))
		} else {
			secondByte, err := reader.ReadByte()

			if err != nil {
				return make([][]int, 0), fmt.Errorf("Error reading second byte: %v", err)
			}

			if secondByte == 0 {
				imageColorIDs = append(imageColorIDs, currentLine)
				currentLine = make([]int, 0)
			} else if secondByte < 64 {
				addIDMultipleTimes(0, 2)
			} else {
				thirdByte, err := reader.ReadByte()

				if err != nil {
					return make([][]int, 0), fmt.Errorf("Error reading third byte: %v", err)
				}

				if secondByte < 128 {
					count := int((secondByte-64))*256 + int(thirdByte)
					addIDMultipleTimes(0, count)

				} else if secondByte < 192 {
					count := int(secondByte) - 128

					addIDMultipleTimes(int(thirdByte), count)

				} else {
					forthByte, err := reader.ReadByte()

					if err != nil {
						return make([][]int, 0), fmt.Errorf("Error reading forth byte: %v", err)
					}

					count := int(secondByte-192)*256 + int(thirdByte)
					addIDMultipleTimes(int(forthByte), count)
				}
			}
		}
	}

	return imageColorIDs, nil
}
