package export

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nandesh-dev/subtle/internal/actions"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logger"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"gorm.io/gorm"
)

func Run() {
	logger.Logger().Log("Export Routine", "Running export routine")
	defer logger.Logger().Log("Export Routine", "Export routine complete")

	var routineEntry database.Routine
	if err := database.Database().Where(database.Routine{Name: "Export"}).FirstOrCreate(&routineEntry, database.Routine{Name: "Export", Description: "Exports subtitles to files.", IsRunning: false}).Error; err != nil {
		logger.Logger().Error("Export Routine", fmt.Errorf("Error getting routine entry from database: %v", err))
		return
	}

	if routineEntry.IsRunning {
		logger.Logger().Error("Export Routine", fmt.Errorf("Media routine is already running"))
		return
	}

	routineEntry.IsRunning = true
	if err := database.Database().Save(routineEntry).Error; err != nil {
		logger.Logger().Error("Export Routine", fmt.Errorf("Error updating routine status in database: %v", err))
		return
	}

	defer func() {
		routineEntry.IsRunning = false
		if err := database.Database().Save(routineEntry).Error; err != nil {
			logger.Logger().Error("Export Routine", fmt.Errorf("Error updating routine status in database: %v", err))
		}
	}()

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
		if !mediaDirectoryConfig.Exporting.Enable {
			continue
		}

		var videoEntries []database.Video
		database.Database().
			Where("directory_path LIKE ?", fmt.Sprintf("%s%%", mediaDirectoryConfig.Path)).
			Preload("Subtitles").
			FindInBatches(&videoEntries, 10, func(tx *gorm.DB, batch int) error {
				for _, videoEntry := range videoEntries {
					for _, subtitleEntry := range videoEntry.Subtitles {
						logger.Logger().Log("Export Routine", fmt.Sprintf("Checking subtitle: %v", subtitleEntry.Title))
						if subtitleEntry.IsExtracted && subtitleEntry.IsFormated && !subtitleEntry.IsExported {
							exportFormat, err := subtitle.ParseFormat(mediaDirectoryConfig.Exporting.Format)
							if err != nil {
								logger.Logger().Error("Export Routine", fmt.Errorf("Error parsing export format in config: %v", err))
								continue
							}

							baseFilepath := filepath.Join(videoEntry.DirectoryPath, strings.Trim(filepath.Base(videoEntry.Filename), filepath.Ext(videoEntry.Filename)))

							if _, err := os.Stat(baseFilepath); err == nil {
								logger.Logger().Error("Export Routine", fmt.Errorf("A subtitle file already exist with the filename: %v", baseFilepath))
								continue
							} else if !os.IsNotExist(err) {
								logger.Logger().Error("Export Routine", fmt.Errorf("Error checking if a subtitle file already exist with the filename: %v", baseFilepath))
								continue
							}

							if err := actions.ExportSubtitle(subtitleEntry.ID, actions.ExportSubtitleConfig{
								Format:       exportFormat,
								BaseFilepath: baseFilepath,
							}); err != nil {
								logger.Logger().Error("Export Routine", fmt.Errorf("Error exporting subtitle: %v", err))
							}
						}
					}
				}
				return nil
			})
	}
}
