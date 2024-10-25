package media

import (
	"context"
	"fmt"
	"path/filepath"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/proto/media"
	"github.com/nandesh-dev/subtle/pkgs/db"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
)

func (s ServiceHandler) GetVideo(ctx context.Context, req *connect.Request[media.GetVideoRequest]) (*connect.Response[media.GetVideoResponse], error) {
	var videoEntry db.Video

	if err := db.DB().Where(&db.Video{DirectoryPath: req.Msg.DirectoryPath, Filename: req.Msg.Name + req.Msg.Extension}).
		Preload("Subtitles").
		Preload("Subtitles.Segments").
		First(&videoEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting video entry: %v", err)
	}

	rawStreams, err := filemanager.NewVideoFile(filepath.Join(req.Msg.DirectoryPath, req.Msg.Name+req.Msg.Extension)).RawStreams()

	if err != nil {
		return nil, fmt.Errorf("Error extracting available raw stream from video: %v", err)
	}

	res := media.GetVideoResponse{
		Subtitles:  make([]*media.Subtitle, len(videoEntry.Subtitles)),
		RawStreams: make([]*media.RawStream, len(*rawStreams)),
	}

	for i, subtitleEntry := range videoEntry.Subtitles {
		res.Subtitles[i] = &media.Subtitle{
			Language:               subtitleEntry.Language,
			ImportIsExternal:       subtitleEntry.ImportIsExternal,
			ImportVideoStreamIndex: int32(subtitleEntry.ImportVideoStreamIndex),
			ExportPath:             subtitleEntry.ExportPath,
		}
	}

	for i, rawStream := range *rawStreams {
		res.RawStreams[i] = &media.RawStream{
			Index:    int32(rawStream.Index()),
			Format:   subtitle.MapFormat(rawStream.Format()),
			Language: rawStream.Language().String(),
			Title:    rawStream.Title(),
		}
	}

	return connect.NewResponse(&res), nil
}
