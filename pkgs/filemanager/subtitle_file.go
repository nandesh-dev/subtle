package filemanager

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
)

type SubtitleFile struct {
	path string
}

func (s *SubtitleFile) Path() string {
	return s.path
}

func (s *SubtitleFile) Basename() string {
	pt := strings.Split(strings.TrimSuffix(filepath.Base(s.path), s.Extension()), ".")

	if len(pt) >= 2 {
		return strings.Join(pt[:len(pt)-1], ".")
	}

	return pt[0]
}

func (s *SubtitleFile) Language() (language.Tag, error) {
	pt := strings.Split(filepath.Base(s.path), ".")

	if len(pt) >= 3 {
		return language.Parse(pt[len(pt)-2])
	}

	return language.English, fmt.Errorf("No language code found")
}

func (s *SubtitleFile) Extension() string {
	return filepath.Ext(s.path)
}
