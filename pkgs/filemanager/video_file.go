package filemanager

import (
	"path/filepath"
	"strings"
)

type VideoFile struct {
	path string
}

func (v *VideoFile) DirectoryPath() string {
	return filepath.Dir(v.path)
}

func (v *VideoFile) Filepath() string {
	return v.path
}

func (v *VideoFile) Filename() string {
	return filepath.Base(v.path)
}

func (v *VideoFile) Extension() string {
	return filepath.Ext(v.path)
}

func (v *VideoFile) Basename() string {
	return strings.TrimSuffix(filepath.Base(v.path), v.Extension())
}
