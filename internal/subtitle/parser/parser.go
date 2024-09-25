package parser

import (
	"fmt"

	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs"
)

func ParseRawSubtitleStream(rawStream subtitle.RawSubtitleStream) (*subtitle.Subtitle, error) {
	switch rawStream.Format {
	case "hdmv_pgs_subtitle":
		return pgs.DecodePGSSubtitle(rawStream)
	}

	return nil, fmt.Errorf("Unsupported subtitle format: %v", rawStream.Format)
}
