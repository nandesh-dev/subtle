package format

import (
	"fmt"

	"github.com/nandesh-dev/subtle/internal/actions"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logger"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"golang.org/x/text/language"
)

func Run() {
	logger.Logger().Log("Format Routine", "Running format routine")
	defer logger.Logger().Log("Format Routine", "Format routine complete")

	var routineEntry database.Routine
	if err := database.Database().Where(database.Routine{Name: "Format"}).FirstOrCreate(&routineEntry, database.Routine{Name: "Format", Description: "Converts the original text / image into final text applying all the formating specified.", IsRunning: false}).Error; err != nil {
		logger.Logger().Error("Format Routine", fmt.Errorf("Error getting routine entry from database: %v", err))
		return
	}

	if routineEntry.IsRunning {
		logger.Logger().Error("Format Routine", fmt.Errorf("Media routine is already running"))
		return
	}

	routineEntry.IsRunning = true
	if err := database.Database().Save(routineEntry).Error; err != nil {
		logger.Logger().Error("Format Routine", fmt.Errorf("Error updating routine status in database: %v", err))
		return
	}

	defer func() {
		routineEntry.IsRunning = false
		if err := database.Database().Save(routineEntry).Error; err != nil {
			logger.Logger().Error("Format Routine", fmt.Errorf("Error updating routine status in database: %v", err))
		}
	}()

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
		if !mediaDirectoryConfig.Formating.Enable {
			continue
		}

		var videoEntries []database.Video
		database.Database().
			Where("directory_path LIKE ?", fmt.Sprintf("%s%%", mediaDirectoryConfig.Path)).
			Preload("Subtitles").
			Find(&videoEntries)

		for _, videoEntry := range videoEntries {
			for _, subtitleEntry := range videoEntry.Subtitles {
				logger.Logger().Log("Format Routine", fmt.Sprintf("Checking subtitle: %v", subtitleEntry.Title))

				if subtitleEntry.IsProcessing {
					continue
				}

				if !subtitleEntry.IsExtracted {
					continue
				}

				if subtitleEntry.IsFormated {
					continue
				}

				if err := database.Database().Where(subtitleEntry).Preload("Segments").Find(&subtitleEntry).Error; err != nil {
					logger.Logger().Error("Format Routine", fmt.Errorf("Error getting subtitle entry from database: %v", err))
					continue
				}

				format, err := subtitle.ParseFormat(subtitleEntry.ImportFormat)
				if err != nil {
					logger.Logger().Error("Format Routine", fmt.Errorf("Error parsing subtitle format: %v", err))
					continue
				}

				lang, err := language.Parse(subtitleEntry.Language)
				if err != nil {
					logger.Logger().Error("Format Routine", fmt.Errorf("Error parsing subtitle language: %v", err))
					continue
				}

				charactorMapping := make([]actions.FormatSubtitleCharactorMapping, 0)
				switch format {
				case subtitle.ASS:
					for _, charactorMappingConfig := range mediaDirectoryConfig.Formating.TextBasedSubtitle.CharactorMappings {
						if charactorMappingConfig.Language == lang {
							for _, mappingConfig := range charactorMappingConfig.Mappings {
								charactorMapping = append(charactorMapping, actions.FormatSubtitleCharactorMapping{
									From: mappingConfig.From,
									To:   mappingConfig.To,
								})
							}
						}
					}
				case subtitle.PGS:
					for _, charactorMappingConfig := range mediaDirectoryConfig.Formating.ImageBasedSubtitle.CharactorMappings {
						if charactorMappingConfig.Language == lang {
							for _, mappingConfig := range charactorMappingConfig.Mappings {
								charactorMapping = append(charactorMapping, actions.FormatSubtitleCharactorMapping{
									From: mappingConfig.From,
									To:   mappingConfig.To,
								})
							}
						}
					}
				default:
					logger.Logger().Error(("Format Routine"), fmt.Errorf("Unsupported format: %v", subtitleEntry.ImportFormat))
					continue
				}

				if err := actions.FormatSubtitle(subtitleEntry.ID, actions.FormatSubtitleConfig{CharactorMappings: charactorMapping}); err != nil {
					logger.Logger().Error("Format Routine", fmt.Errorf("Error formatting subtitle: %v", err))
				}
			}
		}
	}
}
