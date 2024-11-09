package actions

import (
	"fmt"
	"os"

	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logger"
	"github.com/nandesh-dev/subtle/pkgs/srt"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
)

type ExportSubtitleConfig struct {
	Format       subtitle.Format
	BaseFilepath string
}

func ExportSubtitle(id int, config ExportSubtitleConfig) error {
	subtitleEntry := database.Subtitle{
		ID: id,
	}
	if err := database.Database().Where(subtitleEntry).Preload("Segments").First(&subtitleEntry).Error; err != nil {
		return fmt.Errorf("Error getting subtitle entry from database: %v", err)
	}

	subtitleEntry.IsProcessing = true
	if err := database.Database().Save(&subtitleEntry).Error; err != nil {
		return fmt.Errorf("Error updating subtitle processing status: %v", err)
	}

	defer func() {
		subtitleEntry.IsProcessing = false
		if err := database.Database().Save(&subtitleEntry).Error; err != nil {
			logger.Logger().Error("Export Subtitle Action", fmt.Errorf("Error updating subtitle processing status: %v", err))
		}
	}()

	if config.Format != subtitle.SRT {
		return fmt.Errorf("Unsupported export format")
	}

	exportSubtitle := subtitle.NewTextSubtitle()

	for _, segmentEntry := range subtitleEntry.Segments {
		exportSubtitle.AddSegment(*subtitle.NewTextSegment(
			segmentEntry.StartTime,
			segmentEntry.EndTime,
			segmentEntry.Text),
		)
	}

	exportPath := fmt.Sprintf("%s.%s.srt", config.BaseFilepath, subtitleEntry.Language)

	encodedSubtitle := srt.EncodeSubtitle(*exportSubtitle)
	if err := os.WriteFile(exportPath, []byte(encodedSubtitle), 0644); err != nil {
		return fmt.Errorf("Error writting subtitle to file: %v", err)
	}

	subtitleEntry.IsExported = true
	subtitleEntry.ExportPath = exportPath
	subtitleEntry.ExportFormat = subtitle.MapFormat(config.Format)
	if err := database.Database().Save(&subtitleEntry).Error; err != nil {
		return fmt.Errorf("Error updating subtitle entry: %v", err)
	}

	return nil
}
