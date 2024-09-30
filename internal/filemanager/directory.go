package filemanager

import (
	"fmt"

	"github.com/nandesh-dev/subtle/pkgs/warning"
)

type Directory struct {
	path      string
	childrens []Directory
	files     []File
}

func NewDirectory(path string) *Directory {
	return &Directory{
		path:      path,
		childrens: make([]Directory, 0),
		files:     make([]File, 0),
	}
}

func (d *Directory) AddFile(file File) {
	d.files = append(d.files, file)
}

func (d *Directory) AddChild(child Directory) {
	d.childrens = append(d.childrens, child)
}

func (d *Directory) VideoFiles() ([]File, warning.WarningList) {
	videos := make([]File, 0)
	warnings := warning.NewWarningList()

	for _, file := range d.files {
		isVideoFile, err := file.IsVideoFile()

		if err != nil {
			warnings.AddWarning(fmt.Errorf("Error checking is file is of type video: %v", err))
			continue
		}

		if isVideoFile {
			videos = append(videos, file)
		}
	}

	return videos, *warnings
}

func (d *Directory) SubtitleFiles() []File {
	subtitles := make([]File, 0)

	for _, file := range d.files {
		if file.IsSubtitleFile() {
			subtitles = append(subtitles, file)
		}
	}

	return subtitles
}
