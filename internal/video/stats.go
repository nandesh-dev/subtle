package video

import (
	"encoding/json"
	"errors"
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
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
		return nil, errors.New(fmt.Sprintf("Error parsing probe JSON for file %v", v.Path))
	}

	var result map[string]any
	json.Unmarshal([]byte(probeJSON), &result)

	rawStreams, streamsExist := result["streams"].([]any)
	if !streamsExist {
		return nil, errors.New(fmt.Sprintf("Error parsing probe stream JSON for file %v", v.Path))
	}

	streams := make([]Stream, 0)

	for _, rawStream := range rawStreams {
		codecType, codecTypeExist := rawStream.(map[string]any)["codec_type"].(string)
		if !codecTypeExist {
			return nil, errors.New(fmt.Sprintf("Error parsing probe stream codes type JSON for file %v", v.Path))
		}

		if codecType != "subtitle" {
			continue
		}

		codecName, codecNameExist := rawStream.(map[string]any)["codec_name"].(string)
		if !codecNameExist {
			return nil, errors.New(fmt.Sprintf("Error parsing probe stream codec name JSON for file %v", v.Path))
		}

		rawIndex, indexExist := rawStream.(map[string]any)["index"].(float64)
		if !indexExist {
			return nil, errors.New(fmt.Sprintf("Error parsing probe stream index JSON for file %v", v.Path))
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
