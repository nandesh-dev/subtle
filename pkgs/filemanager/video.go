package filemanager

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nandesh-dev/subtle/pkgs/warning"
	"golang.org/x/text/language"
)

type VideoFile struct {
	path      string
	subtitles []SubtitleFile
}

func NewVideoFile(path string) *VideoFile {
	return &VideoFile{
		path: path,
	}
}

func (v *VideoFile) DirectoryPath() string {
	return filepath.Dir(v.path)
}

func (v *VideoFile) Filepath() string {
	return v.path
}

func (v *VideoFile) Extension() string {
	return filepath.Ext(v.path)
}

func (v *VideoFile) Basename() string {
	return strings.TrimSuffix(filepath.Base(v.path), v.Extension())
}

func (v *VideoFile) HasSubtitleLanguage(tag language.Tag) (bool, warning.WarningList) {
	warnings := warning.NewWarningList()
	for _, subtitleFile := range v.SubtitleFiles() {
		subtitleLanguageTag, err := language.Parse(subtitleFile.LanguageCode())
		if err != nil {
			warnings.AddWarning(fmt.Errorf("Invalid language code: %v; %v", subtitleFile.LanguageCode(), err))
			continue
		}

		if subtitleLanguageTag == tag {
			return true, *warnings
		}
	}

	return false, *warnings
}

func (v *VideoFile) AddSubtitleFile(file SubtitleFile) {
	v.subtitles = append(v.subtitles, file)
}

func (v *VideoFile) SubtitleFiles() []SubtitleFile {
	return v.subtitles
}
