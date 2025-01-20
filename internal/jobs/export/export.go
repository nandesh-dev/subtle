package export

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	cue_schema "github.com/nandesh-dev/subtle/generated/ent/cue"
	cue_content_segment_schema "github.com/nandesh-dev/subtle/generated/ent/cuecontentsegment"
	subtitle_schema "github.com/nandesh-dev/subtle/generated/ent/subtitle"
	video_schema "github.com/nandesh-dev/subtle/generated/ent/video"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/subtitle/srt"
)

func Run(logger *slog.Logger, conf *config.Config, db *ent.Client) {
	c, err := conf.Read()
	if err != nil {
		logger.Error("cannot read config", "err", err)
		return
	}

	for _, mediaDirectoryConfig := range c.MediaDirectories {
		if !mediaDirectoryConfig.Exporting.Enable {
			continue
		}

		exportFormat, err := subtitle.ParseFormat(mediaDirectoryConfig.Exporting.Format)
		if err != nil {
			logger.Error("cannot parse export format from config", "err", err)
			continue
		}

		if exportFormat != subtitle.SRT {
			logger.Error("unsupported export format")
			continue
		}

		videoEntries, err := db.Video.Query().Where(video_schema.FilepathHasPrefix(mediaDirectoryConfig.Path)).All(context.Background())
		if err != nil {
			logger.Error("cannot get videos from database", "err", err)
			continue
		}

		for _, videoEntry := range videoEntries {
			logger := logger.With("video_filepath", videoEntry.Filepath)

			subtitleEntries, err := videoEntry.QuerySubtitles().
				Where(subtitle_schema.IsProcessingEQ(false), subtitle_schema.StageEQ(subtitle_schema.StageFormated)).
				All(context.Background())
			if err != nil {
				logger.Error("cannot get subtitle from database", "err", err)
			}

			for _, subtitleEntry := range subtitleEntries {
				logger := logger.With("subtitle_title", subtitleEntry.Title)

				logger.Info("exporting subtitle")

				exportFilepath := fmt.Sprintf(
					"%s.%s.%s",
					strings.TrimSuffix(videoEntry.Filepath, filepath.Ext(videoEntry.Filepath)),
					subtitleEntry.Language,
					subtitle.MapFormat(exportFormat),
				)

				if _, err := os.Stat(exportFilepath); err == nil {
					logger.Error("subtitle file already exist with same filepath")
					continue
				} else if !os.IsNotExist(err) {
					logger.Error("cannot check if subtitle file already exist", "err", err)
					continue
				}

				tx, err := db.Tx(context.Background())
				if err != nil {
					logger.Error(logging.DatabaseTransactionCreateError, "err", err)
					continue
				}

				if err := db.Subtitle.UpdateOne(subtitleEntry).SetIsProcessing(true).Exec(context.Background()); err != nil {
					logger.Error("cannot update subtitle processing status in database", "err", err)
					continue
				}

				defer func() {
					if err := db.Subtitle.UpdateOneID(subtitleEntry.ID).SetIsProcessing(false).Exec(context.Background()); err != nil {
						logger.Error("cannot update subtitle processing status in database", "err", err)
					}
				}()

				if err := func() error {
					cueEntries, err := tx.Cue.Query().
						Where(cue_schema.HasSubtitleWith(subtitle_schema.ID(subtitleEntry.ID))).
						Order(cue_schema.ByTimestampStart()).
						All(context.Background())
					if err != nil {
						logger.Error("cannot get subtitle cues from database", "err", err)
						return err
					}

					sub := subtitle.Subtitle{}
					for _, cueEntry := range cueEntries {
						content := make([]subtitle.CueContentSegment, 0)

						contentEntries, err := tx.CueContentSegment.Query().
							Where(cue_content_segment_schema.HasCueWith(cue_schema.IDEQ(cueEntry.ID))).
							Order(cue_content_segment_schema.ByPosition()).
							All(context.Background())
						if err != nil {
							logger.Error("cannot get subtitle cue content from database", "err", err)
							return err
						}

						for _, cueContentEntry := range contentEntries {
							content = append(content, subtitle.CueContentSegment{
								Text: cueContentEntry.Text,
							})
						}

						sub.Cues = append(sub.Cues, subtitle.Cue{
							Timestamp: subtitle.CueTimestamp{
								Start: cueEntry.TimestampStart,
								End:   cueEntry.TimestampEnd,
							},
							Content: content,
						})
					}

					file, err := os.OpenFile(exportFilepath, os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						logger.Error("cannot open export file for writing", "err", err)
						return err
					}

          defer file.Close()

					if err := sub.Write(srt.NewWriter(file)); err != nil {
						logger.Error("cannot write subtitle to file", "err", err)
						return err
					}

					if err := tx.Subtitle.UpdateOne(subtitleEntry).
						SetStage(subtitle_schema.StageExported).
						Exec(context.Background()); err != nil {
						logger.Error("cannot update subtitle to exported", "err", err)
						return err
					}

					logger.Info("subtitle exported")

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
