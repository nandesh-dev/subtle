package media

import (
	"fmt"
	"os"

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

	var routineEntry database.Routine
	if err := database.Database().Where(database.Routine{Name: "Media"}).FirstOrCreate(&routineEntry, database.Routine{Name: "Media", Description: "Scans the media directory for new video files and extract raw subtitle streams from it.", IsRunning: false}).Error; err != nil {
		logger.Logger().Error("Media Routine", fmt.Errorf("Error getting routine entry from database: %v", err))
		return
	}

	if routineEntry.IsRunning {
		logger.Logger().Error("Media Routine", fmt.Errorf("Media routine is already running"))
		return
	}

	routineEntry.IsRunning = true
	if err := database.Database().Save(routineEntry).Error; err != nil {
		logger.Logger().Error("Media Routine", fmt.Errorf("Error updating routine status in database: %v", err))
		return
	}

	defer func() {
		routineEntry.IsRunning = false
		if err := database.Database().Save(routineEntry).Error; err != nil {
			logger.Logger().Error("Media Routine", fmt.Errorf("Error updating routine status in database: %v", err))
		}
	}()

	for _, mediaDirectoryConfig := range config.Config().MediaDirectories {
		logger.Logger().Log("Media Routine", fmt.Sprintf("Reading watch directory: %v", mediaDirectoryConfig.Path))

		watchDirectory, _, err := filemanager.ReadDirectory(mediaDirectoryConfig.Path, true)
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

				for i, subtitleEntry := range videoEntry.Subtitles {
					if subtitleEntry.IsExported {
						if _, err := os.Stat(subtitleEntry.ExportPath); err != nil {
							if os.IsNotExist(err) {
								videoEntry.Subtitles[i].IsExported = false
								videoEntry.Subtitles[i].ExportPath = ""
								videoEntry.Subtitles[i].ExportFormat = ""

								logger.Logger().Log("Media Routine", fmt.Sprintf("Subtitle missing exported file, reverted to unexported state: %v", subtitleEntry.Title))
							} else {
								logger.Logger().Error("Media Routine", fmt.Errorf("Error checking if subtitle file exist: %v", err))
							}
						}
					}
				}

				if err := database.Database().Save(videoEntry).Error; err != nil {
					logger.Logger().Error("Media Routine", fmt.Errorf("Error saving video entry: %v", err))
				}
			}
		}
	}
}
