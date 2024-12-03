package filemanager

import (
	"os"
	"path/filepath"
)

func ReadDirectory(path string) (*Directory, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	directory := &Directory{
		Path:          path,
		ChildrenPaths: make([]string, 0),
		Videos:        make([]VideoFile, 0),
		Subtitles:     make([]SubtitleFile, 0),
	}

	for _, entry := range files {
		entrypath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			directory.ChildrenPaths = append(directory.ChildrenPaths, entrypath)
		}

		if IsSubtitleFile(entrypath) {
			directory.Subtitles = append(directory.Subtitles, SubtitleFile{
				path: entrypath,
			})
			continue
		}

		isVideoFile, err := IsVideoFile(entrypath)

		if err != nil {
			continue
		}

		if isVideoFile {
			directory.Videos = append(directory.Videos, VideoFile{
				path: entrypath,
			})
		}
	}

	return directory, nil
}
