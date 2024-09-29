package decoder

import (
	"fmt"

	"github.com/nandesh-dev/subtle/internal/ass"
	"github.com/nandesh-dev/subtle/internal/pgs"
	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/warnings"
)

func DecodeRawSubtitleStream(rawStream subtitle.RawStream) (*subtitle.Stream, error, *warnings.WarningList) {
	switch rawStream.Format {
	case subtitle.PGS:
		return pgs.DecodeSubtitle(rawStream)
	case subtitle.ASS:
		return ass.DecodeSubtitle(rawStream)
	}

	return nil, fmt.Errorf("Unsupported subtitle format: %v", rawStream.Format), nil
}
