package tesseract

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log/slog"

	"github.com/otiai10/gosseract/v2"
	"golang.org/x/text/language"
)

type TesseractClient struct {
	GosseractClient *gosseract.Client
}

func NewClient() *TesseractClient {
	return &TesseractClient{
		GosseractClient: gosseract.NewClient(),
	}
}

func (c *TesseractClient) Close() {
	if err := c.GosseractClient.Close(); err != nil {
		slog.Warn("Error closing gosseract client: %v", err)
	}
}

func (c *TesseractClient) ExtractTextFromImage(img image.Image, lang language.Tag) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("Error encoding image to png: %v", err)
	}

	if err := c.GosseractClient.SetImageFromBytes(buf.Bytes()); err != nil {
		return "", fmt.Errorf("Error scanning image: %v", err)
	}

	text, err := c.GosseractClient.Text()
	if err != nil {
		return "", fmt.Errorf("Error getting text from gosseract: %v", err)
	}

	return text, nil
}
