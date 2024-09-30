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

	directory := NewDirectory(path)

	for _, entry := range files {
		if entry.IsDir() {
			child, err := ReadDirectory(filepath.Join(path, entry.Name()))
			if err != nil {
				return nil, err
			}

			directory.AddChild(*child)
		}

		file := NewFile(filepath.Join(path, entry.Name()))
		directory.AddFile(*file)
	}

	return directory, nil
}
