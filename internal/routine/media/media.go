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
					Raw(
						"SELECT COUNT(*) FROM videos WHERE directory_path = ? AND filename = ?;",
						video.DirectoryPath(), video.Filename(),
					).Scan(&count).Error; err != nil {
					logger.Error("cannot cound videos in database", "err", err)
					continue
				}

				if count > 0 {
					continue
				}

				logger.Info("processing video", slog.Group("info", "path", video.Filepath()))

				var videoId int64
				if err := database.Database().
					Raw(
						"INSERT INTO videos (directory_path, filename, is_processing) VALUES (?, ?, ?) RETURNING id;",
						video.DirectoryPath(), video.Filename(), true,
					).Scan(&videoId).Error; err != nil {
					logger.Error("cannot insert video into database", "err", err)
					continue
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
					lang := rawStream.Language()

					if err := database.Database().
						Exec(
							"INSERT INTO subtitles (video_id, language, import_is_external, import_format, import_video_stream_index) VALUES (?, ?, ?, ?, ?);",
							videoId, lang.String(), false, subtitle.MapFormat(rawStream.Format()), rawStream.Index(),
						).Error; err != nil {
						logger.Error("cannot insert subtitle into database", "err", err)
					}
				}

				if err := database.Database().
					Exec(
						"UPDATE videos SET is_processing=? WHERE id = ?;",
						false, videoId,
					).Error; err != nil {
					logger.Error("cannot update video processing status in database", "err", err)
				}
			}
		}
	}
}
