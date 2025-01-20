package media

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nandesh-dev/subtle/generated/ent"
	subtitle_schema "github.com/nandesh-dev/subtle/generated/ent/subtitle"
	video_schema "github.com/nandesh-dev/subtle/generated/ent/video"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"golang.org/x/text/language/display"
)

func Run(logger *slog.Logger, conf *config.Config, db *ent.Client) {
	c, err := conf.Read()
	if err != nil {
		logger.Error("cannot read config", "err", err)
		return
	}

	// Loop through all media directories and scan them for new video files
	for _, mediaDirectoryConfig := range c.MediaDirectories {
		directoryPathStack := []string{mediaDirectoryConfig.Path}

		for len(directoryPathStack) > 0 {
			path := directoryPathStack[len(directoryPathStack)-1]
			directoryPathStack = directoryPathStack[:len(directoryPathStack)-1]

			logger.Info("reading directory", slog.Group("info", "path", path))
			directory, err := filemanager.ReadDirectory(path)
			if err != nil {
				logger.Error("cannot read directory", "err", err)
				continue
			}

			directoryPathStack = append(directoryPathStack, directory.ChildrenPaths...)

			for _, video := range directory.Videos {
				tx, err := db.Tx(context.Background())
				if err != nil {
					logger.Error("cannot create database transaction", "err", err)
					continue
				}

				if err := func() error {
					// Skip if the video already exist in the database
					if existingVideoEntryCount, err := tx.Video.
						Query().
						Where(video_schema.Filepath(video.Filepath())).
						Count(context.Background()); err != nil {
						logger.Error("cannot count number of extracted videos", "err", err)
						return err
					} else if existingVideoEntryCount > 0 {
						return nil
					}

					logger.Info("processing video", slog.Group("info", "path", video.Filepath()))

					videoEntry, err := tx.Video.
						Create().
						SetFilepath(video.Filepath()).
						Save(context.Background())
					if err != nil {
						logger.Error("cannot save video to database", "err", err)
						return err
					}

					rawStreams, err := video.RawStreams()
					if err != nil {
						logger.Error("cannot get raw stream from video", "err", err)
						return err
					}

					for _, rawStream := range *rawStreams {
						title := rawStream.Title()

						if title == "" {
							title = fmt.Sprintf("%v#%v", display.Self.Name(rawStream.Language()), rawStream.Index())
						}

						if err := tx.Subtitle.
							Create().
							SetTitle(title).
							SetLanguage(rawStream.Language().String()).
							SetStage(subtitle_schema.StageDetected).
							SetImportFormat(subtitle.MapFormat(rawStream.Format())).
							SetImportVideoStreamIndex(int32(rawStream.Index())).
							AddVideo(videoEntry).
							Exec(context.Background()); err != nil {
							logger.Error("cannot save subtitle to database", "err", err)
							return err
						}
					}

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
