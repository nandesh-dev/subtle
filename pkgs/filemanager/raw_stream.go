package filemanager

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/warning"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/text/language"
)

type RawStream struct {
	filepath string
	index    int
	format   subtitle.Format
	language language.Tag
	title    string
}

func (s *RawStream) Filepath() string {
	return s.filepath
}

func (s *RawStream) Index() int {
	return s.index
}

func (s *RawStream) Format() subtitle.Format {
	return s.format
}

func (s *RawStream) Language() language.Tag {
	return s.language
}

func (s *RawStream) Title() string {
	return s.title
}

func (v *VideoFile) RawStreams() (*[]RawStream, error) {
	rawResultData, err := ffmpeg.Probe(v.path)
	warnings := warning.NewWarningList()

	if err != nil {
		return nil, fmt.Errorf("Failed to probe file stats: %v", err)
	}

	if !json.Valid([]byte(rawResultData)) {
		return nil, fmt.Errorf("Invalid probe JSON data: %v", rawResultData)
	}

	var result map[string]any
	json.Unmarshal([]byte(rawResultData), &result)

	rawStreamsData, streamsExist := result["streams"].([]any)
	if !streamsExist {
		return nil, fmt.Errorf("Missing streams in probe JSON: %v", rawResultData)
	}

	rawStreams := make([]RawStream, 0)

	for _, rawStreamData := range rawStreamsData {
		codecType, codecTypeExist := rawStreamData.(map[string]any)["codec_type"].(string)
		if !codecTypeExist {
			return nil, fmt.Errorf("Missing codec type in probe JSON: %v", rawStreamData)
		}

		if codecType != "subtitle" {
			continue
		}

		codecName, codecNameExist := rawStreamData.(map[string]any)["codec_name"].(string)
		if !codecNameExist {
			return nil, fmt.Errorf("Missing codec name in probe JSON: %v", rawStreamData)
		}

		format, err := mapCodecName(codecName)
		if err != nil {
			return nil, err
		}

		rawIndex, indexExist := rawStreamData.(map[string]any)["index"].(float64)
		if !indexExist {
			return nil, fmt.Errorf("Missing index in probe JSON: %v", rawStreamData)
		}

		lang := language.English
		title := ""

		tags, tagsExist := rawStreamData.(map[string]any)["tags"].(map[string]any)
		if tagsExist {
			rawLanguage, langaugeExist := tags["language"].(string)

			if langaugeExist {
				langTag, err := language.Parse(rawLanguage)
				if err != nil {
					warnings.AddWarning(fmt.Errorf("Invalid language in probe JSON: %v; %v", rawLanguage, rawStreamData))
				} else {
					lang = langTag
				}
			}

			if rawTitle, titleExist := tags["title"].(string); titleExist {
				title = strings.TrimSpace(rawTitle)
			}
		}

		rawStreams = append(rawStreams, RawStream{
			filepath: v.path,
			index:    int(rawIndex),
			format:   format,
			language: lang,
			title:    title,
		})
	}

	return &rawStreams, nil
}

func mapCodecName(cN string) (subtitle.Format, error) {
	switch cN {
	case "hdmv_pgs_subtitle":
		return subtitle.PGS, nil
	case "ass":
		return subtitle.ASS, nil
	case "subrip":
		return subtitle.SRT, nil
	}

	return subtitle.ASS, fmt.Errorf("Unsupported or invalid codec name: %v", cN)
}
