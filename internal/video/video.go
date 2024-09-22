package video

import (
	"encoding/json"
	"errors"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoFile struct {
	Path   string
	Format VideoFileFileFormat
}

type VideoFileFileFormat int

const (
	MP4 VideoFileFileFormat = iota
	MKV
	AVI
	MOV
)

type Stream struct {
	Index  int
	Format string
}

type VideoStats struct {
	Streams []Stream
}

func (v *VideoFile) Stats() (*VideoStats, error) {
	probeJSON, err := ffmpeg.Probe(v.Path)
	if err != nil {
		return nil, err
	}

	if !json.Valid([]byte(probeJSON)) {
		return nil, errors.New("Error parsing probe JSON")
	}

	var result map[string]any
	json.Unmarshal([]byte(probeJSON), &result)

	rawStreams, streamsExist := result["streams"].([]any)
	if !streamsExist {
		return nil, errors.New("Error parsing probe stream JSON")
	}

	streams := make([]Stream, 0)

	for _, rawStream := range rawStreams {
		codecType, codecTypeExist := rawStream.(map[string]any)["codec_type"].(string)
		if !codecTypeExist {
			return nil, errors.New("Error parsing probe stream codes type JSON")
		}

		if codecType != "subtitle" {
			continue
		}

		codecName, codecNameExist := rawStream.(map[string]any)["codec_name"].(string)
		if !codecNameExist {
			return nil, errors.New("Error parsing probe stream codec name JSON")
		}

		rawIndex, indexExist := rawStream.(map[string]any)["index"].(float64)
		if !indexExist {
			return nil, errors.New("Error parsing probe stream index JSON")
		}

		stream := Stream{
			Index:  int(rawIndex),
			Format: codecName,
		}

		streams = append(streams, stream)
	}

	return &VideoStats{
		Streams: streams,
	}, nil
}
