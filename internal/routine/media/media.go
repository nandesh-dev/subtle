package media

import (
	"fmt"
	"log/slog"

	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/logging"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"golang.org/x/text/language/display"
)

func Run() {
	logger := logging.NewRoutineLogger("media")

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
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
				var count int64
				if err := database.Database().
					Model(&database.Video{}).
					Where(database.Video{
						DirectoryPath: video.DirectoryPath(),
						Filename:      video.Filename(),
					}).Count(&count).Error; err != nil {
					logger.Error("cannot cound videos in database", "err", err)
					continue
				}

				if count > 0 {
					continue
				}

				logger.Info("processing video", slog.Group("info", "path", video.Filepath()))

				videoEntry := database.Video{
					DirectoryPath: video.DirectoryPath(),
					Filename:      video.Filename(),

					Subtitles: make([]database.Subtitle, 0),
				}

				rawStreams, err := video.RawStreams()
				if err != nil {
					logger.Error("cannot get raw stream from video", "err", err)
					continue
				}

				for _, rawStream := range *rawStreams {
					title := rawStream.Title()

					if title == "" {
						title = fmt.Sprintf("%v#%v", display.Self.Name(rawStream.Language()), rawStream.Index())
					}

					subtitleEntry := database.Subtitle{
						Language: rawStream.Language().String(),

						ImportIsExternal:       false,
						ImportFormat:           subtitle.MapFormat(rawStream.Format()),
						ImportVideoStreamIndex: rawStream.Index(),
					}

					videoEntry.Subtitles = append(videoEntry.Subtitles, subtitleEntry)
				}

				logger.Info("adding video to database", slog.Group("info", "filepath", video.Filepath()))
				if err := database.Database().Create(&videoEntry).Error; err != nil {
					logger.Error("cannot save video to database", "err", err)
				}
			}
		}
	}
}
