package actions

import (
	"bytes"
	"fmt"
	"image/png"
	"path/filepath"

	"github.com/nandesh-dev/subtle/pkgs/ass"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/logger"
	"github.com/nandesh-dev/subtle/pkgs/pgs"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/tesseract"
	"golang.org/x/text/language"
)

func ExtractSubtitle(subtitleEntry database.Subtitle) {
	logger.Logger().Log("Extract Subtitle Action", fmt.Sprintf("Extracting Subtitle; ID: %v", subtitleEntry.ID))

	if subtitleEntry.IsProcessing {
		logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Subtitle is already being processed: ID = %v", subtitleEntry.ID))
		return
	}

	subtitleEntry.IsProcessing = true
	if err := database.Database().Save(subtitleEntry).Error; err != nil {
		logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error updating subtitle processing status: %v", err))
		return
	}

	defer func() {
		subtitleEntry.IsProcessing = false
		if err := database.Database().Save(subtitleEntry).Error; err != nil {
			logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error updating subtitle processing status: %v", err))
		}
	}()

	format, err := subtitle.ParseFormat(subtitleEntry.ImportFormat)
	if err != nil {
		logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error parsing subtitle format: %v", err))
		return
	}

	lang, err := language.Parse(subtitleEntry.Language)
	if err != nil {
		logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error parsing subtitle language: %v", err))
		return

	}

	videoEntry := database.Video{
		ID: subtitleEntry.VideoID,
	}

	if err := database.Database().Where(videoEntry).First(&videoEntry).Error; err != nil {
		logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error getting video entry: %v", err))
		return
	}

	rawStream := filemanager.NewRawStream(
		filepath.Join(videoEntry.DirectoryPath, videoEntry.Filename),
		subtitleEntry.ImportVideoStreamIndex,
		format,
		lang,
		subtitleEntry.Title,
	)

	var sub subtitle.Subtitle

	switch format {
	case subtitle.ASS:
		s, _, err := ass.ExtractFromRawStream(*rawStream)
		if err != nil {
			logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error extracting ASS subtitle: %v", err))
			return
		}

		sub = s
	case subtitle.PGS:
		s, _, err := pgs.ExtractFromRawStream(*rawStream)
		if err != nil {
			logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error extracting ASS subtitle: %v", err))
			return
		}

		sub = s

	default:
		logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Unsupported subtitle format: %v", subtitleEntry.ImportFormat))
		return
	}

	subtitleEntry.Segments = make([]database.Segment, 0)

	switch sub := sub.(type) {
	case subtitle.TextSubtitle:
		for _, segment := range sub.Segments() {
			segmentEntry := database.Segment{
				StartTime:    segment.Start(),
				EndTime:      segment.End(),
				OriginalText: segment.Text(),
			}

			subtitleEntry.Segments = append(subtitleEntry.Segments, segmentEntry)
		}
	case subtitle.ImageSubtitle:
		tesseractClient := tesseract.NewClient()
		defer tesseractClient.Close()

		for _, segment := range sub.Segments() {
			imageDataBuffer := new(bytes.Buffer)
			if err := png.Encode(imageDataBuffer, segment.Image()); err != nil {
				logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error encoding image to png: %v", err))
				continue
			}

			segmentEntry := database.Segment{
				StartTime:     segment.Start(),
				EndTime:       segment.End(),
				OriginalImage: imageDataBuffer.Bytes(),
			}

			subtitleEntry.Segments = append(subtitleEntry.Segments, segmentEntry)
		}
	}

	subtitleEntry.IsExtracted = true

	if err := database.Database().Save(subtitleEntry).Error; err != nil {
		logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error saving subtitle to database: %v", err))
	}
}
