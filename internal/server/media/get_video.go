package media

import (
	"context"
	"fmt"
	"path/filepath"

	"connectrpc.com/connect"
	media_proto "github.com/nandesh-dev/subtle/generated/proto/media"
	subtitle_proto "github.com/nandesh-dev/subtle/generated/proto/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/db"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
)

func (s ServiceHandler) GetVideo(ctx context.Context, req *connect.Request[media_proto.GetVideoRequest]) (*connect.Response[media_proto.GetVideoResponse], error) {
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

	res := media_proto.GetVideoResponse{
		Subtitles:  make([]*subtitle_proto.Subtitle, len(videoEntry.Subtitles)),
		RawStreams: make([]*media_proto.RawStream, len(*rawStreams)),
	}

	for i, subtitleEntry := range videoEntry.Subtitles {
		res.Subtitles[i] = &subtitle_proto.Subtitle{
			Id:               int32(subtitleEntry.ID),
			Title:            "Hello World",
			Language:         subtitleEntry.Language,
			ImportIsExternal: subtitleEntry.ImportIsExternal,
			ExportPath:       subtitleEntry.ExportPath,
		}
	}

	for i, rawStream := range *rawStreams {
		res.RawStreams[i] = &media_proto.RawStream{
			Index:    int32(rawStream.Index()),
			Format:   subtitle.MapFormat(rawStream.Format()),
			Language: rawStream.Language().String(),
			Title:    rawStream.Title(),
		}
	}

	return connect.NewResponse(&res), nil
}
