package ocr

import (
	"bytes"
	"image"
	"image/png"

	"github.com/otiai10/gosseract/v2"
)

type OCRClient struct {
	GosseractClient *gosseract.Client
}

func NewOCRClient() *OCRClient {
	client := gosseract.NewClient()

	return &OCRClient{
		GosseractClient: client,
	}
}

func (c *OCRClient) Close() {
	c.GosseractClient.Close()
}

func (c *OCRClient) ScanImage(img image.Image) (string, error) {
	var buf bytes.Buffer
	png.Encode(&buf, img)

	c.GosseractClient.SetImageFromBytes(buf.Bytes())

	text, err := c.GosseractClient.Text()

	if err != nil {
		return "", err
	}

	return text, nil
}
