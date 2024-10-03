package filemanager

import (
	"path/filepath"
	"strings"
)

type SubtitleFile struct {
	path string
}

func NewSubtitleFile(path string) *SubtitleFile {
	return &SubtitleFile{
		path: path,
	}
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

func (s *SubtitleFile) LanguageCode() string {
	pt := strings.Split(filepath.Base(s.path), ".")

	if len(pt) >= 3 {
		return pt[len(pt)-2]
	}

	return ""
}

func (s *SubtitleFile) Extension() string {
	return filepath.Ext(s.path)
}
