package filemanager

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type File struct {
	path string
}

func NewFile(path string) *File {
	return &File{
		path: path,
	}
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Extension() string {
	return filepath.Ext(f.path)
}

func (f *File) Basename() string {
	return strings.TrimSuffix(filepath.Base(f.path), f.Extension())
}

func (f *File) IsVideoFile() (bool, error) {
	rawProbeResult, err := ffmpeg.Probe(f.path)
	if err != nil {
		return false, fmt.Errorf("Failed to probe file stats: %v", err)
	}

	if !json.Valid([]byte(rawProbeResult)) {
		return false, fmt.Errorf("Invalid probe JSON data: %v", rawProbeResult)
	}

	var result map[string]any
	json.Unmarshal([]byte(rawProbeResult), &result)

	rawStreams, streamsExist := result["streams"].([]any)
	if !streamsExist {
		return false, fmt.Errorf("Missing streams in probe JSON: %v", rawProbeResult)
	}

	for _, rawStream := range rawStreams {
		codecType, codecTypeExist := rawStream.(map[string]any)["codec_type"].(string)
		if !codecTypeExist {
			return false, fmt.Errorf("Missing codec type in probe JSON: %v", rawStream)
		}

		if codecType == "video" {
			return true, nil
		}
	}

	return false, nil
}

func (f *File) IsSubtitleFile() bool {
	return slices.Contains([]string{}, f.Extension())
}
