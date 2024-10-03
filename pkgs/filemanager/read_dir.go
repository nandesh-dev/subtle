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

	directory := NewDirectory(path)
	videos := make([]*VideoFile, 0)
	subtitles := make([]*SubtitleFile, 0)

	for _, entry := range files {
		entrypath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			child, w, err := ReadDirectory(entrypath)
			warnings.Append(w)

			if err != nil {
				return nil, *warnings, err
			}

			directory.AddChild(*child)
		}

		if IsSubtitleFile(entrypath) {
			subtitles = append(subtitles, NewSubtitleFile(entrypath))
			continue
		}

		isVideoFile, err := IsVideoFile(entrypath)

		if err != nil {
			warnings.AddWarning(fmt.Errorf("Error checking if file is video: %v", err))
			continue
		}

		if isVideoFile {
			videos = append(videos, NewVideoFile(entrypath))
		}
	}

	extraSubtitles := make([]*SubtitleFile, 0)

	for _, subtitle := range subtitles {
		found := false
		for _, video := range videos {
			if subtitle.Basename() == video.Basename() {
				found = true
				video.AddSubtitleFile(*subtitle)
				break
			}
		}

		if !found {
			extraSubtitles = append(extraSubtitles, subtitle)
		}
	}

	for _, video := range videos {
		directory.AddVideoFile(*video)
	}

	for _, subtitle := range extraSubtitles {
		directory.AddExtraSubtitleFile(*subtitle)
	}

	return directory, *warnings, nil
}
