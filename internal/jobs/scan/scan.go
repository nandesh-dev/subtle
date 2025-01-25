package scan

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/nandesh-dev/subtle/generated/ent"
	"github.com/nandesh-dev/subtle/generated/ent/subtitleschema"
	"github.com/nandesh-dev/subtle/generated/ent/videoschema"
	"github.com/nandesh-dev/subtle/pkgs/configuration"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/logging"
)

func Run(ctx context.Context, logger *slog.Logger, configFile *configuration.File, db *ent.Client) {
	config, err := configFile.Read()
	if err != nil {
		logger.Error("cannot read config", "err", err)
		return
	}

	for _, scanningGroup := range config.Job.Scanning {
		directoryPathStack := []string{scanningGroup.DirectoryPath}

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
				handleVideoProcessing(video, ctx, logger, db)
			}
		}
	}
}

func handleVideoProcessing(video filemanager.VideoFile, ctx context.Context, logger *slog.Logger, db *ent.Client) {
	tx, err := db.Tx(ctx)
	if err != nil {
		logger.Error("cannot create database transaction", "err", err)
		return
	}

	if err := func() error {
		if existingVideoEntryCount, err := tx.VideoSchema.
			Query().
			Where(videoschema.Filepath(video.Filepath())).
			Count(ctx); err != nil {
			logger.Error("cannot count number of extracted videos", "err", err)
			return err
		} else if existingVideoEntryCount > 0 {
			return nil
		}

		logger.Info("processing video", slog.Group("info", "path", video.Filepath()))

		videoEntry, err := tx.VideoSchema.
			Create().
			SetFilepath(video.Filepath()).
			Save(ctx)
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
			if strings.TrimSpace(title) == "" {
				title = fmt.Sprintf("%v#%v", rawStream.Language().String(), rawStream.Index())
			}

			if err := tx.SubtitleSchema.
				Create().
				SetTitle(title).
				SetLanguage(rawStream.Language()).
				SetStage(subtitleschema.StageDetected).
				SetImportFormat(rawStream.Format()).
				SetImportVideoStreamIndex(rawStream.Index()).
				AddVideo(videoEntry).
				Exec(ctx); err != nil {
				logger.Error("cannot save subtitle to database", "err", err)
				return err
			}
		}

		return nil
	}(); err != nil {
    logger.Warn("rolling back")

		if err := tx.Rollback(); err != nil {
			logger.Error(logging.DatabaseTransactionRollbackError, "err", err)
		}
	} else {
		if err := tx.Commit(); err != nil {
			logger.Error(logging.DatabaseTransactionCommitError, "err", err)
		}
	}
}
