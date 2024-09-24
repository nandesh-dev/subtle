package srt

import (
	"fmt"
	"log"

	"github.com/nandesh-dev/subtle/internal/ocr"
	"github.com/nandesh-dev/subtle/internal/subtitle"
)

func EncodeSRTSubtitles(st subtitle.Subtitle) {
	images := st.ImageStream

	ocrClient := ocr.NewOCRClient()
	defer ocrClient.Close()

	for _, segment := range images {
		for _, img := range segment.Images {
			text, err := ocrClient.ScanImage(img)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(text + "\n")
		}
	}
}
