package format

import (
	"bytes"
	"context"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	"github.com/nandesh-dev/subtle/generated/ent/segment"
	subtitle_schema "github.com/nandesh-dev/subtle/generated/ent/subtitle"
	video_schema "github.com/nandesh-dev/subtle/generated/ent/video"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/tesseract"
	"golang.org/x/text/language"
)

func Run(conf *config.Config, db *ent.Client) {
	logger := logging.NewRoutineLogger("format")

	c, err := conf.Read()
	if err != nil {
		logger.Error("cannot read config", "err", err)
		return
	}

	for _, mediaDirectoryConfig := range c.MediaDirectories {
		if !mediaDirectoryConfig.Formating.Enable {
			continue
		}

		videoEntries, err := db.Video.Query().
			Where(video_schema.FilepathHasPrefix(mediaDirectoryConfig.Path)).
			All(context.Background())
		if err != nil {
			logger.Error("cannot get videos from database", "err", err)
			continue
		}

		for _, videoEntry := range videoEntries {
			logger := logger.With("video_filepath", videoEntry.Filepath)

			subtitleEntries, err := videoEntry.QuerySubtitles().
				Where(subtitle_schema.Processing(false), subtitle_schema.Extracted(true), subtitle_schema.Formated(false), subtitle_schema.Exported(false)).
				All(context.Background())
			if err != nil {
				logger.Error("cannot get video subtitles from database", "err", err)
				continue
			}

			for _, subtitleEntry := range subtitleEntries {
				logger := logger.With("subtitle_title", subtitleEntry.Title)

				if err := db.Subtitle.UpdateOne(subtitleEntry).SetProcessing(true).Exec(context.Background()); err != nil {
					logger.Error("cannot update subtitle processing status in database", "err", err)
					continue
				}

				defer func() {
					if err := db.Subtitle.UpdateOneID(subtitleEntry.ID).SetProcessing(false).Exec(context.Background()); err != nil {
						logger.Error("cannot update subtitle processing status in database", "err", err)
					}
				}()

				tx, err := db.Tx(context.Background())

				if err != nil {
					logger.Error(logging.DatabaseTransactionCreateError, "err", err)
					continue
				}

				if err := func() error {
					logger.Info("formating subtitle")

					format, err := subtitle.ParseFormat(subtitleEntry.ImportFormat)
					if err != nil {
						logger.Error("cannot parse subtitle format", "err", err)
						return err
					}

					lang, err := language.Parse(subtitleEntry.Language)
					if err != nil {
						logger.Error("cannot parse subtitle language", "err", err)
						return err
					}

					imageBasedSubtitle := false
					if format == subtitle.PGS {
						imageBasedSubtitle = true
					}

					segmentEntries, err := tx.Segment.Query().
						Where(segment.HasSubtitleWith(subtitle_schema.ID(subtitleEntry.ID))).
						All(context.Background())
					if err != nil {
						logger.Error("cannot get subtitle segments from database", "err", err)
						return err
					}

					tesseractClient := tesseract.NewClient()
					defer tesseractClient.Close()
					for _, segmentEntry := range segmentEntries {
						extractedText := segmentEntry.OriginalText
						charactorMappings := mediaDirectoryConfig.Formating.TextBasedSubtitle.CharactorMappings

						if imageBasedSubtitle {
							text, err := tesseractClient.ExtractTextFromPNGImage(*bytes.NewBuffer(segmentEntry.OriginalImage), lang)
							if err != nil {
								logger.Error("cannot extract from image", "err", err)
								return err
							}
							extractedText = text
							charactorMappings = mediaDirectoryConfig.Formating.ImageBasedSubtitle.CharactorMappings
						}

						formatedText := extractedText

						for _, charactorMapping := range charactorMappings {
							if lang == charactorMapping.Language {
								for _, mapping := range charactorMapping.Mappings {
									formatedText = strings.ReplaceAll(formatedText, mapping.From, mapping.To)
								}
							}
						}

						if err := tx.Segment.UpdateOne(segmentEntry).SetText(formatedText).Exec(context.Background()); err != nil {
							logger.Error("cannot update segment in database", "err", err)
							return err
						}
					}

					if err := tx.Subtitle.UpdateOne(subtitleEntry).SetFormated(true).Exec(context.Background()); err != nil {
						logger.Error("cannot update subtitle formated status", "err", err)
						return err
					}

					logger.Info("subtitle formated")
					return nil
				}(); err != nil {
					if err := tx.Rollback(); err != nil {
						logger.Error(logging.DatabaseTransactionRollbackError, "err", err)
					}
				} else {
					if err := tx.Commit(); err != nil {
						logger.Error(logging.DatabaseTransactionCommitError, "err", err)
					}
				}
			}
		}
	}
}
