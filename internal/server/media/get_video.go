package media

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"

	"connectrpc.com/connect"
	media_proto "github.com/nandesh-dev/subtle/generated/proto/media"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
)

func (s ServiceHandler) GetVideo(ctx context.Context, req *connect.Request[media_proto.GetVideoRequest]) (*connect.Response[media_proto.GetVideoResponse], error) {
	var videoEntry database.Video

	if err := database.Database().Where(&database.Video{ID: int(req.Msg.Id)}).
		Preload("Subtitles").
		First(&videoEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting video entry: %v", err)
	}

	video := filemanager.NewVideoFile(filepath.Join(videoEntry.DirectoryPath, videoEntry.Filename))

	res := media_proto.GetVideoResponse{
		Id:                 req.Msg.Id,
		DirectoryPath:      video.DirectoryPath(),
		BaseName:           video.Basename(),
		Extension:          video.Extension(),
		SubtitleIds:        make([]int32, 0),
		IsProcessing:       false,
		ExtractedLanguages: make([]string, 0),
	}

	for _, subtitleEntry := range videoEntry.Subtitles {
		res.SubtitleIds = append(res.SubtitleIds, int32(subtitleEntry.ID))

		if subtitleEntry.IsProcessing {
			res.IsProcessing = true
		}

		if subtitleEntry.IsExported {
			if !slices.Contains(res.ExtractedLanguages, subtitleEntry.Language) {
				res.ExtractedLanguages = append(res.ExtractedLanguages, subtitleEntry.Language)
			}
		}
	}

	return connect.NewResponse(&res), nil
}
