package extract

import (
	"bytes"
	"fmt"
	"image/png"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/nandesh-dev/subtle/pkgs/ass"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/pgs"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/tesseract"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

func Run() {
	logger := logging.NewRoutineLogger("extract")

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
		if !mediaDirectoryConfig.Extraction.Enable {
			continue
		}

		var videoEntries []database.Video
		database.Database().
			Where("directory_path LIKE ?", fmt.Sprintf("%s%%", mediaDirectoryConfig.Path)).
			Preload("Subtitles").
			FindInBatches(&videoEntries, 10, func(tx *gorm.DB, batch int) error {
				for _, videoEntry := range videoEntries {
					alreadyExtracted := false
					for _, subtitleEntry := range videoEntry.Subtitles {
						if subtitleEntry.IsExtracted {
							alreadyExtracted = true
							break
						}
					}

					if alreadyExtracted {
						continue
					}

					logger.Info("looking for suitable subtitle to extract")

					bestScore := -1
					var bestSubtitleEntry *database.Subtitle

					for _, subtitleEntry := range videoEntry.Subtitles {
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

					if bestSubtitleEntry == nil {
						continue
					}

					subtitleEntry := bestSubtitleEntry

					subtitleEntry.IsProcessing = true
					if err := database.Database().Save(&subtitleEntry).Error; err != nil {
						logger.Error("cannot save subtitle to database", "err", err)
					}

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
							logger.Error("cannot extract subtitle", "err", err, slog.Group("info", "format", "ass"))
							continue
						}

						sub = *s
					case subtitle.PGS:
						s, _, err := pgs.ExtractFromRawStream(*rawStream)
						if err != nil {
							logger.Error("cannot extract subtitle", "err", err, slog.Group("info", "format", "pgs"))
							continue
						}

						sub = *s

					default:
						logger.Error("unsupported / invalid subtitle format", slog.Group("info", "format", format))
						continue
					}

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
								logger.Error("cannot encode image", "err", err)
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

					subtitleEntry.IsProcessing = true
					subtitleEntry.IsExtracted = true
					if err := database.Database().Save(&subtitleEntry).Error; err != nil {
						logger.Error("cannot save subtitle to database", "err", err)
					}
				}

				return nil
			})
	}
}
