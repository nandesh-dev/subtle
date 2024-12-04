package extract

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/nandesh-dev/subtle/internal/actions"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"golang.org/x/text/language"
)

func Run() {
	logger := logging.NewRoutineLogger("extract")

	var routineEntry database.Routine
	if err := database.Database().
		Where(database.Routine{Name: "Extract"}).
		FirstOrCreate(
			&routineEntry,
			database.Routine{
				Name:        "Extract",
				Description: "Converts the raw subtitle streams into usable text / images based subtitles.",
				IsRunning:   false,
			},
		).Error; err != nil {
		logger.Error("cannot get routine from database", "err", err)
		return
	}

	if routineEntry.IsRunning {
		logger.Error("already running")
		return
	}

	routineEntry.IsRunning = true
	if err := database.Database().Save(routineEntry).Error; err != nil {
		logger.Error("cannot update subtitle status in database", "err", err)
		return
	}

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
		if !mediaDirectoryConfig.Extraction.Enable {
			continue
		}

		var videoEntries []database.Video
		database.Database().
			Where("directory_path LIKE ?", fmt.Sprintf("%s%%", mediaDirectoryConfig.Path)).
			Preload("Subtitles").
			Find(&videoEntries)

		for _, videoEntry := range videoEntries {
			bestScore := -1
			var bestSubtitleEntry *database.Subtitle

			for _, subtitleEntry := range videoEntry.Subtitles {
				if subtitleEntry.IsExtracted {
					break
				}

				logger.Info("evaluating subtitle", slog.Group("info", "title", subtitleEntry.Title, "format", subtitleEntry.ImportFormat))

				format, err := subtitle.ParseFormat(subtitleEntry.ImportFormat)
				if err != nil {
					logger.Error("cannot parse subtitle format", "err", err)
					continue
				}

				lang, err := language.Parse(subtitleEntry.Language)
				if err != nil {
					logger.Error("cannot parse subtitle language", "err", err)
					continue
				}

				containsRequiredLanguage := false

				switch format {
				case subtitle.ASS:
					if !mediaDirectoryConfig.Extraction.Formats.ASS.Enable {
						continue
					}

					for _, languageTag := range mediaDirectoryConfig.Extraction.Formats.ASS.Languages {
						if lang == languageTag {
							containsRequiredLanguage = true
						}
					}
				case subtitle.PGS:
					if !mediaDirectoryConfig.Extraction.Formats.PGS.Enable {
						continue
					}

					for _, languageTag := range mediaDirectoryConfig.Extraction.Formats.PGS.Languages {
						if lang == languageTag {
							containsRequiredLanguage = true
						}
					}
				}

				if !containsRequiredLanguage {
					continue
				}

				score := 0

				for _, rawStreamTitleKeyword := range mediaDirectoryConfig.Extraction.RawStreamTitleKeywords {
					if strings.Contains(subtitleEntry.Title, rawStreamTitleKeyword) {
						score++
					}
				}

				if score > bestScore {
					bestScore = score
					bestSubtitleEntry = &subtitleEntry
				}
			}

			if bestSubtitleEntry != nil {
				actions.ExtractSubtitle(bestSubtitleEntry.ID)
			}
		}
	}

	routineEntry.IsRunning = false
	if err := database.Database().Save(routineEntry).Error; err != nil {
		logger.Error("cannot update routine status in database", "err", err)
	}
}
