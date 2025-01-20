package format

import (
	"context"
	"log/slog"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	cue_schema "github.com/nandesh-dev/subtle/generated/ent/cue"
	cue_content_segment_schema "github.com/nandesh-dev/subtle/generated/ent/cuecontentsegment"
	subtitle_schema "github.com/nandesh-dev/subtle/generated/ent/subtitle"
	video_schema "github.com/nandesh-dev/subtle/generated/ent/video"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"golang.org/x/text/language"
)

func Run(logger *slog.Logger, conf *config.Config, db *ent.Client) {
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
				Where(subtitle_schema.IsProcessingEQ(false), subtitle_schema.StageEQ(subtitle_schema.StageExtracted)).
				All(context.Background())
			if err != nil {
				logger.Error("cannot get video subtitles from database", "err", err)
				continue
			}

			for _, subtitleEntry := range subtitleEntries {
				logger := logger.With("subtitle_title", subtitleEntry.Title)

				if err := db.Subtitle.UpdateOne(subtitleEntry).SetIsProcessing(true).Exec(context.Background()); err != nil {
					logger.Error("cannot update subtitle processing status in database", "err", err)
					continue
				}

				defer func() {
					if err := db.Subtitle.UpdateOneID(subtitleEntry.ID).SetIsProcessing(false).Exec(context.Background()); err != nil {
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

					lang, err := language.Parse(subtitleEntry.Language)
					if err != nil {
						logger.Error("cannot parse subtitle language", "err", err)
						return err
					}

					cueEntries, err := tx.Cue.Query().
						Where(cue_schema.HasSubtitleWith(subtitle_schema.IDEQ(subtitleEntry.ID))).
						All(context.Background())
					if err != nil {
						logger.Error("cannot get subtitle segments from database", "err", err)
						return err
					}

					for _, cueEntry := range cueEntries {
						segmentEntries, err := cueEntry.QueryCueContentSegments().
							Order(cue_content_segment_schema.ByPosition()).
							All(context.Background())
						if err != nil {
							logger.Error("cannot query cue content segment from database", "err", err)
							return err
						}

						for _, segmentEntry := range segmentEntries {
							formatedText := segmentEntry.Text
							for _, charactorMapping := range mediaDirectoryConfig.Formating.TextBasedSubtitle.CharactorMappings {
								if lang == charactorMapping.Language {
									for _, mapping := range charactorMapping.Mappings {
										formatedText = strings.ReplaceAll(formatedText, mapping.From, mapping.To)
									}
								}
							}

							if err := segmentEntry.Update().SetText(formatedText).Exec(context.Background()); err != nil {
								logger.Error("cannot save formated cue segment to database", "err", err)
								return err
							}
						}
					}

					if err := tx.Subtitle.UpdateOne(subtitleEntry).
						SetStage(subtitle_schema.StageFormated).
						Exec(context.Background()); err != nil {
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
