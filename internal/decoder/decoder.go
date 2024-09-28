package decoder

import (
	"fmt"

	"github.com/nandesh-dev/subtle/internal/pgs"
	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/warnings"
)

func DecodeRawSubtitleStream(rawStream subtitle.RawStream) (*subtitle.Subtitle, error, *warnings.WarningList) {
	switch rawStream.Format {
	case "hdmv_pgs_subtitle":
		return pgs.DecodePGSSubtitle(rawStream)
	}

	return nil, fmt.Errorf("Unsupported subtitle format: %v", rawStream.Format), nil
}
