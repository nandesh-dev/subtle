package media

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"path/filepath"

	"connectrpc.com/connect"
	media_proto "github.com/nandesh-dev/subtle/generated/proto/media"
	"github.com/nandesh-dev/subtle/pkgs/ass"
	"github.com/nandesh-dev/subtle/pkgs/db"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/pgs"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/tesseract"
)

func (s ServiceHandler) ExtractRawStream(ctx context.Context, req *connect.Request[media_proto.ExtractRawStreamRequest]) (*connect.Response[media_proto.ExtractRawStreamResponse], error) {
	var videoEntry db.Video

	if err := db.DB().Where(&db.Video{ID: int(req.Msg.VideoId)}).First(&videoEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting video entry from database: %v", err)
	}

	video := filemanager.NewVideoFile(filepath.Join(videoEntry.DirectoryPath, videoEntry.Filename))

	rawStreams, err := video.RawStreams()
	if err != nil {
		return nil, fmt.Errorf("Error getting raw streams from video: %v", err)
	}

	for _, rawStream := range *rawStreams {
		if rawStream.Index() == int(req.Msg.RawStreamIndex) {
			var sub subtitle.Subtitle

			switch rawStream.Format() {
			case subtitle.ASS:
				assSubtitle, _, err := ass.DecodeSubtitle(rawStream)
				if err != nil {
					return nil, fmt.Errorf("Error decoding ass subtitle: %v", err)
				}
				sub = assSubtitle
			case subtitle.PGS:
				pgsSubtitle, _, err := pgs.DecodeSubtitle(rawStream)
				if err != nil {
					return nil, fmt.Errorf("Error decoding pgs subtitle: %v", err)
				}
				sub = pgsSubtitle
			default:
				return nil, fmt.Errorf("Invalid or unsupported subtitle codec")
			}

			subtitleEntry := db.Subtitle{
				VideoID:  videoEntry.ID,
				Title:    req.Msg.Title,
				Language: rawStream.Language().String(),
				Segments: make([]db.Segment, 0),
			}

			switch sub := sub.(type) {
			case subtitle.TextSubtitle:
				for _, segment := range sub.Segments() {
					segmentEntry := db.Segment{
						StartTime:    segment.Start(),
						EndTime:      segment.End(),
						Text:         segment.Text(),
						OriginalText: segment.Text(),
					}

					subtitleEntry.Segments = append(subtitleEntry.Segments, segmentEntry)
				}
			case subtitle.ImageSubtitle:
				tesseractClient := tesseract.NewClient()
				defer tesseractClient.Close()

				for _, segment := range sub.Segments() {
					imageDataBuffer := new(bytes.Buffer)
					if err := png.Encode(imageDataBuffer, segment.Image()); err != nil {
						continue
					}

					text, err := tesseractClient.ExtractTextFromPNGImage(*imageDataBuffer, rawStream.Language())
					if err != nil {
						continue
					}

					segmentEntry := db.Segment{
						StartTime:     segment.Start(),
						EndTime:       segment.End(),
						Text:          text,
						OriginalText:  text,
						OriginalImage: imageDataBuffer.Bytes(),
					}

					subtitleEntry.Segments = append(subtitleEntry.Segments, segmentEntry)
				}
			}

			db.DB().Save(&subtitleEntry)
			res := media_proto.ExtractRawStreamResponse{}
			return connect.NewResponse(&res), nil
		}
	}

	return nil, fmt.Errorf("Invalid raw stream index")
}
