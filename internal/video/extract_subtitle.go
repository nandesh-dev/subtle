package video

import (
	"bytes"
	"fmt"

	"github.com/nandesh-dev/subtle/internal/subtitle"
	"github.com/nandesh-dev/subtle/internal/subtitle/pgs"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func (v *VideoFile) ExtractSubtitle(stream Stream) (*subtitle.Subtitle, error) {
	var ouputBuf bytes.Buffer

	ffmpeg.LogCompiledCommand = false
	err := ffmpeg.Input(v.Path).
		Output("pipe:", ffmpeg.KwArgs{"map": fmt.Sprintf("0:%v", stream.Index), "c:s": "copy", "f": "sup"}).
		WithOutput(&ouputBuf).
		Run()

	if err != nil {
		return nil, fmt.Errorf("Error extracting subtitles: %v", err)
	}

	subtitle, err := pgs.DecodePGSSubtitle(ouputBuf.Bytes())

	if err != nil {
		return nil, fmt.Errorf("Error decoding pgs subtitle: %v", err)
	}

	return subtitle, nil
}
