package actions

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logger"
	"github.com/nandesh-dev/subtle/pkgs/tesseract"
	"golang.org/x/text/language"
)

type FormatSubtitleCharactorMapping struct {
	From string
	To   string
}

type FormatSubtitleConfig struct {
	CharactorMappings []FormatSubtitleCharactorMapping
}

func FormatSubtitle(id int, config FormatSubtitleConfig) error {
	logger.Logger().Log("Format Subtitle Action", fmt.Sprintf("Formatting Subtitle; ID: %v", id))

	subtitleEntry := database.Subtitle{ID: id}
	if err := database.Database().Where(subtitleEntry).Preload("Segments").First(&subtitleEntry).Error; err != nil {
		return fmt.Errorf("Error getting subtitle entry from database: %v", err)
	}

	if subtitleEntry.IsFormated {
		return fmt.Errorf("Subtitle is already formated")
	}

	subtitleEntry.IsProcessing = true
	if err := database.Database().Save(subtitleEntry).Error; err != nil {
		return fmt.Errorf("Error updating subtitle processing status: %v", err)
	}

	defer func() {
		subtitleEntry.IsProcessing = false
		if err := database.Database().Save(subtitleEntry).Error; err != nil {
			logger.Logger().Error("Extract Subtitle Action", fmt.Errorf("Error updating subtitle processing status: %v", err))
		}
	}()

	lang, err := language.Parse(subtitleEntry.Language)
	if err != nil {
		return fmt.Errorf("Error parsing subtitle language: %v", err)
	}

	tesseractClient := tesseract.NewClient()
	defer tesseractClient.Close()
	for i, segmentEntry := range subtitleEntry.Segments {
		if strings.TrimSpace(segmentEntry.OriginalText) != "" {
			text := segmentEntry.OriginalText

			for _, charactorMapping := range config.CharactorMappings {
				text = strings.ReplaceAll(text, charactorMapping.From, charactorMapping.To)
			}

			segmentEntry.Text = text
			subtitleEntry.Segments[i].Text = text
		}

		if len(segmentEntry.OriginalImage) > 0 {
			text, err := tesseractClient.ExtractTextFromPNGImage(*bytes.NewBuffer(segmentEntry.OriginalImage), lang)
			if err != nil {
				return fmt.Errorf("Error extracting text from image using tesseract: %v", err)
			}

			for _, charactorMapping := range config.CharactorMappings {
				text = strings.ReplaceAll(text, charactorMapping.From, charactorMapping.To)
			}

			segmentEntry.Text = text
			subtitleEntry.Segments[i].Text = text
		}
	}

	subtitleEntry.IsFormated = true

	if err := database.Database().Save(&subtitleEntry).Error; err != nil {
		return fmt.Errorf("Error saving subtitle to database: %v", err)
	}

	return nil
}
