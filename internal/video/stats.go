package video

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nandesh-dev/subtle/internal/subtitle"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/text/language"
)

type Stats struct {
	RawStreams []subtitle.RawStream
}

func (v *File) Stats() (*Stats, error) {
	probeJSON, err := ffmpeg.Probe(v.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed to probe file stats: %v", err)
	}

	if !json.Valid([]byte(probeJSON)) {
		return nil, fmt.Errorf("Invalid probe JSON data: %v", v.Path)
	}

	var result map[string]any
	json.Unmarshal([]byte(probeJSON), &result)

	rawStreams, streamsExist := result["streams"].([]any)
	if !streamsExist {
		return nil, fmt.Errorf("Missing streams in probe JSON: %v", v.Path)
	}

	streams := make([]subtitle.RawStream, 0)

	for _, rawStream := range rawStreams {
		codecType, codecTypeExist := rawStream.(map[string]any)["codec_type"].(string)
		if !codecTypeExist {
			return nil, fmt.Errorf("Missing codec type in probe JSON: %v", v.Path)
		}

		if codecType != "subtitle" {
			continue
		}

		codecName, codecNameExist := rawStream.(map[string]any)["codec_name"].(string)
		if !codecNameExist {
			return nil, fmt.Errorf("Missing codec name in probe JSON: %v", v.Path)
		}

		format, err := mapCodecName(codecName)
		if err != nil {
			return nil, err
		}

		rawIndex, indexExist := rawStream.(map[string]any)["index"].(float64)
		if !indexExist {
			return nil, fmt.Errorf("Missing index in probe JSON: %v", v.Path)
		}

		lang := language.English

		tags, tagsExist := rawStream.(map[string]any)["tags"].(map[string]any)
		if tagsExist {
			rawLanguage, langaugeExist := tags["language"].(string)

			if langaugeExist {
				langTag, err := language.Parse(rawLanguage)

				if err != nil {
					slog.Warn("Invalid language in probe JSON: %v; %v", rawLanguage, v.Path)
				} else {
					lang = langTag
				}
			}
		}

		stream := subtitle.RawStream{
			Index:         int(rawIndex),
			Format:        format,
			Language:      lang,
			VideoFilePath: v.Path,
		}

		streams = append(streams, stream)
	}

	return &Stats{
		RawStreams: streams,
	}, nil
}

func mapCodecName(cN string) (subtitle.Format, error) {
	switch cN {
	case "hdmv_pgs_subtitle":
		return subtitle.PGS, nil
	case "ass":
		return subtitle.ASS, nil
	}

	return subtitle.ASS, fmt.Errorf("Unsupported or invalid codec name: %v", cN)
}
