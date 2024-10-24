package filemanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nandesh-dev/subtle/pkgs/warning"
)

func ReadDirectory(path string) (*Directory, warning.WarningList, error) {
	warnings := warning.NewWarningList()
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, *warnings, err
	}

	directory := &Directory{
		path:      path,
		children:  make([]Directory, 0),
		videos:    make([]VideoFile, 0),
		subtitles: make([]SubtitleFile, 0),
	}

	for _, entry := range files {
		entrypath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			child, w, err := ReadDirectory(entrypath)
			warnings.Append(w)

			if err != nil {
				return nil, *warnings, err
			}

			directory.children = append(directory.children, *child)
		}

		if IsSubtitleFile(entrypath) {
			directory.subtitles = append(directory.subtitles, SubtitleFile{
				path: entrypath,
			})
			continue
		}

		isVideoFile, err := IsVideoFile(entrypath)

		if err != nil {
			warnings.AddWarning(fmt.Errorf("Error checking if file is video: %v", err))
			continue
		}

		if isVideoFile {
			directory.videos = append(directory.videos, VideoFile{
				path: entrypath,
			})
		}
	}

	return directory, *warnings, nil
}
