package media

import (
	"fmt"

	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/logger"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"golang.org/x/text/language/display"
)

func Run() {
	logger.Logger().Log("Media Routine", "Running media routine")
	defer logger.Logger().Log("Media Routine", "Media routine complete")

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
		logger.Logger().Log("Media Routine", fmt.Sprintf("Reading watch directory: %v", mediaDirectoryConfig.Path))

		watchDirectory, _, err := filemanager.ReadDirectory(mediaDirectoryConfig.Path)
		if err != nil {
			logger.Logger().Error("Media Routine", fmt.Errorf("Error reading watch directory: %v", err))
			return
		}

		stack := []filemanager.Directory{*watchDirectory}

		for len(stack) > 0 {
			currentDirectory := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, currentDirectory.Children()...)

			for _, video := range currentDirectory.VideoFiles() {
				logger.Logger().Log("Media Routine", fmt.Sprintf("Syncing video file: %v", video.Filepath()))

				videoEntry := database.Video{
					DirectoryPath: video.DirectoryPath(),
					Filename:      video.Filename(),
				}

				if err := database.Database().
					Where(videoEntry).
					Preload("Subtitles").
					FirstOrCreate(&videoEntry).
					Error; err != nil {
					logger.Logger().Error("Media Routine", fmt.Errorf("Error creating video entry: %v", err))
					return
				}

				rawStreams, err := video.RawStreams()
				if err != nil {
					logger.Logger().Error("Media Routine", fmt.Errorf("Error getting raw stream from video: %v", err))
					return
				}

				for _, rawStream := range *rawStreams {
					entryExist := false
					for _, subtitleEntry := range videoEntry.Subtitles {
						if !subtitleEntry.ImportIsExternal && subtitleEntry.ImportVideoStreamIndex == rawStream.Index() {
							entryExist = true
							break
						}
					}

					if !entryExist {
						logger.Logger().Log("Media Routine", fmt.Sprintf("Syncing subtitle indexed: %v; of format: %v; from video :%v", rawStream.Index(), subtitle.MapFormat(rawStream.Format()), video.Filepath()))

						subtitleEntry := database.Subtitle{
							Language: rawStream.Language().String(),

							ImportIsExternal:       false,
							ImportFormat:           subtitle.MapFormat(rawStream.Format()),
							ImportVideoStreamIndex: rawStream.Index(),
						}

						if rawStream.Title() != "" {
							subtitleEntry.Title = rawStream.Title()
						} else {
							subtitleEntry.Title = fmt.Sprintf("%v#%v", display.Self.Name(rawStream.Language()), rawStream.Index())
						}
						videoEntry.Subtitles = append(videoEntry.Subtitles, subtitleEntry)
					}
				}

				if err := database.Database().Save(videoEntry).Error; err != nil {
					logger.Logger().Error("Media Routine", fmt.Errorf("Error saving video entry: %v", err))
				}
			}
		}
	}
}
