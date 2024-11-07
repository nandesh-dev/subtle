package extract

import (
	"fmt"
	"strings"

	"github.com/nandesh-dev/subtle/internal/actions"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/logger"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"golang.org/x/text/language"
)

func Run() {
	logger.Logger().Log("Extract Routine", "Running extract routine")
	defer logger.Logger().Log("Extract Routine", "Extract routine complete")

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
		var videoEntries []database.Video
		database.Database().
			Where("directory_path LIKE ?", fmt.Sprintf("%s%%", mediaDirectoryConfig.Path)).
			Preload("Subtitles").
			Find(&videoEntries)

		for _, videoEntry := range videoEntries {
			bestScore := -1
			var bestSubtitleEntry *database.Subtitle

			for _, subtitleEntry := range videoEntry.Subtitles {
				logger.Logger().Log("Media Routine", fmt.Sprintf("Checking subtitle: %v", subtitleEntry.Title))
				if subtitleEntry.IsExtracted {
					break
				}

				logger.Logger().Log("Media Routine", fmt.Sprintf("Extracting subtitle: %v", subtitleEntry.Title))

				format, err := subtitle.ParseFormat(subtitleEntry.ImportFormat)
				if err != nil {
					logger.Logger().Error("Extract Routine", fmt.Errorf("Error parsing subtitle format: %v", err))
					continue
				}

				lang, err := language.Parse(subtitleEntry.Language)
				if err != nil {
					logger.Logger().Error("Extract Routine", fmt.Errorf("Error parsing subtitle language: %v", err))
					continue
				}

				containsRequiredLanguage := false

				switch format {
				case subtitle.ASS:
					for _, languageTag := range mediaDirectoryConfig.Extraction.Formats.ASS.Languages {
						if lang == languageTag {
							containsRequiredLanguage = true
						}
					}
				case subtitle.PGS:
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
				actions.ExtractSubtitle(*bestSubtitleEntry)
			}
		}
	}
}
