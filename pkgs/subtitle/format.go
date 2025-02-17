package subtitle

import (
	"database/sql/driver"
	"fmt"

	"github.com/nandesh-dev/subtle/generated/proto/web"
	"gopkg.in/yaml.v3"
)

type code int

const (
	pgs code = iota
	ass
	srt
)

type Format struct {
	yaml.Marshaler
	yaml.Unmarshaler
	code code
}

func (f Format) FileExt() string {
	switch f.code {
	case pgs:
		return "sup"
	case ass:
		return "ass"
	case srt:
		return "srt"
	}

	return ""
}

func (f Format) String() string {
	switch f.code {
	case pgs:
		return "PGS"
	case ass:
		return "ASS"
	case srt:
		return "SRT"
	}

	return ""
}

func ParseFormat(str string) (*Format, error) {
	switch str {
	case "PGS":
		return &PGS, nil
	case "ASS":
		return &ASS, nil
	case "SRT":
		return &SRT, nil
	}

	return nil, fmt.Errorf("invalid format \"%s\"", str)
}

var (
	PGS Format = Format{code: pgs}
	ASS Format = Format{code: ass}
	SRT Format = Format{code: srt}
)

func (f Format) FFMpegString() string {
	switch f.code {
	case pgs:
		return "sup"
	case ass:
		return "ass"
	case srt:
		return "srt"
	}

	return ""
}

func (f Format) MarshalYAML() (interface{}, error) {
	return f.String(), nil
}

func (f *Format) UnmarshalYAML(value *yaml.Node) error {
	var rawFormat string
	if err := value.Decode(&rawFormat); err != nil {
		return err
	}

	format, err := ParseFormat(rawFormat)
	if err != nil {
		return err
	}

	f.code = format.code

	return nil
}

func (f Format) Value() (driver.Value, error) {
	return f.String(), nil
}

func (f *Format) Scan(src any) error {
	switch v := src.(type) {
	case string:
		format, err := ParseFormat(v)
		if err != nil {
			return err
		}

		f.code = format.code
	default:
		return fmt.Errorf("unsupposed type: %t", src)
	}

	return nil
}

func (f *Format) ToProto() web.SubtitleImportFormat {
  switch f.code {
  case srt:
    return web.SubtitleImportFormat_SRT
  case ass:
    return web.SubtitleImportFormat_ASS
  case pgs:
    return web.SubtitleImportFormat_PGS
  }

  return web.SubtitleImportFormat_SRT
}
