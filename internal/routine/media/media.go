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

	var routineEntry database.Routine
	if err := database.Database().
		Where(database.Routine{Name: "Media"}).
		FirstOrCreate(
			&routineEntry,
			database.Routine{
				Name:        "Media",
				Description: "Scans the media directory for new video files and extract raw subtitle streams from it.",
				IsRunning:   false,
			},
		).Error; err != nil {
		logger.Error("cannot get routine from database", "err", err)
		return
	}

	if routineEntry.IsRunning {
		logger.Error("routine already running")
		return
	}

	routineEntry.IsRunning = true
	if err := database.Database().Save(routineEntry).Error; err != nil {
		logger.Error("cannot update routine status in database", "err", err)
		return
	}

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
		logger.Info("reading watch directory", slog.Group("info", "path", mediaDirectoryConfig.Path))

		watchDirectory, _, err := filemanager.ReadDirectory(mediaDirectoryConfig.Path, true)
		if err != nil {
			logger.Error("cannot react watch directory", "err", err)
			return
		}

		stack := []filemanager.Directory{*watchDirectory}

		for len(stack) > 0 {
			currentDirectory := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, currentDirectory.Children()...)

			for _, video := range currentDirectory.VideoFiles() {
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

				videoEntry := database.Video{
					DirectoryPath: video.DirectoryPath(),
					Filename:      video.Filename(),
				}

				logger.Info("adding video to database", slog.Group("info", "filepath", video.Filepath()))
				if err := database.Database().Create(&videoEntry).Error; err != nil {
					logger.Error("cannot save video to database", "err", err)
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
						VideoID:  videoEntry.ID,
						Language: rawStream.Language().String(),

						ImportIsExternal:       false,
						ImportFormat:           subtitle.MapFormat(rawStream.Format()),
						ImportVideoStreamIndex: rawStream.Index(),
					}

					logger.Info("adding raw stream subtitle to database", slog.Group("info", "title", title, "video_id", videoEntry.ID))
					if err := database.Database().Create(&subtitleEntry).Error; err != nil {
						logger.Error("cannot save subtitle to database", "err", err)
					}
				}
			}
		}
	}
}
