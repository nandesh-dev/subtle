package filemanager

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func IsVideoFile(path string) (bool, error) {
	rawProbeResult, err := ffmpeg.Probe(path)
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

func IsSubtitleFile(path string) bool {
	return slices.Contains([]string{".srt"}, filepath.Ext(path))
}
