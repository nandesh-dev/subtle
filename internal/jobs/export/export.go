package export

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	"github.com/nandesh-dev/subtle/generated/ent/segment"
	subtitle_schema "github.com/nandesh-dev/subtle/generated/ent/subtitle"
	video_schema "github.com/nandesh-dev/subtle/generated/ent/video"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/srt"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
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
				Where(subtitle_schema.Processing(false), subtitle_schema.Extracted(true), subtitle_schema.Formated(true), subtitle_schema.Exported(false)).
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

				if err := db.Subtitle.UpdateOne(subtitleEntry).SetProcessing(true).Exec(context.Background()); err != nil {
					logger.Error("cannot update subtitle processing status in database", "err", err)
					continue
				}

				defer func() {
					if err := db.Subtitle.UpdateOneID(subtitleEntry.ID).SetProcessing(false).Exec(context.Background()); err != nil {
						logger.Error("cannot update subtitle processing status in database", "err", err)
					}
				}()

				if err := func() error {
					segmentEntries, err := tx.Segment.Query().Where(segment.HasSubtitleWith(subtitle_schema.ID(subtitleEntry.ID))).All(context.Background())
					if err != nil {
						logger.Error("cannot get subtitle segments from database", "err", err)
						return err
					}

					exportSubtitle := subtitle.NewTextSubtitle()

					for _, segmentEntry := range segmentEntries {
						exportSubtitle.AddSegment(*subtitle.NewTextSegment(
							segmentEntry.StartTime,
							segmentEntry.EndTime,
							segmentEntry.Text),
						)
					}

					encodedSubtitle := srt.EncodeSubtitle(*exportSubtitle)

					if err := os.WriteFile(exportFilepath, []byte(encodedSubtitle), 0644); err != nil {
						logger.Error("cannot write subtitle to file", "err", err)
						return err
					}

					if err := tx.Subtitle.UpdateOne(subtitleEntry).SetExported(true).Exec(context.Background()); err != nil {
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
