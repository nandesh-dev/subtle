package subtitle

import (
	"encoding/json"
	"fmt"

	"github.com/nandesh-dev/subtle/pkgs/warning"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/text/language"
)

type Format int

const (
	PGS Format = iota
	ASS
)

type RawStream struct {
	filepath string
	index    int
	format   Format
	language language.Tag
}

func NewRawStream(filepath string) *RawStream {
	return &RawStream{
		filepath: filepath,
	}
}

func (s *RawStream) Filepath() string {
	return s.filepath
}

func (s *RawStream) SetIndex(index int) {
	s.index = index
}

func (s *RawStream) Index() int {
	return s.index
}

func (s *RawStream) SetFormat(format Format) {
	s.format = format
}

func (s *RawStream) Format() Format {
	return s.format
}

func (s *RawStream) SetLanguage(lang language.Tag) {
	s.language = lang
}

func (s *RawStream) Language() language.Tag {
	return s.language
}

type file interface {
	Path() string
}

func ExtractRawStreams(file file) ([]RawStream, error) {
	rawResultData, err := ffmpeg.Probe(file.Path())
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

		rawStream := NewRawStream(file.Path())

		codecName, codecNameExist := rawStreamData.(map[string]any)["codec_name"].(string)
		if !codecNameExist {
			return nil, fmt.Errorf("Missing codec name in probe JSON: %v", rawStreamData)
		}

		format, err := mapCodecName(codecName)
		if err != nil {
			return nil, err
		}
		rawStream.SetFormat(format)

		rawIndex, indexExist := rawStreamData.(map[string]any)["index"].(float64)
		if !indexExist {
			return nil, fmt.Errorf("Missing index in probe JSON: %v", rawStreamData)
		}

		rawStream.SetIndex(int(rawIndex))

		lang := language.English

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
		}

		rawStream.SetLanguage(lang)

		rawStreams = append(rawStreams, *rawStream)
	}

	return rawStreams, nil
}

func mapCodecName(cN string) (Format, error) {
	switch cN {
	case "hdmv_pgs_subtitle":
		return PGS, nil
	case "ass":
		return ASS, nil
	}

	return ASS, fmt.Errorf("Unsupported or invalid codec name: %v", cN)
}
