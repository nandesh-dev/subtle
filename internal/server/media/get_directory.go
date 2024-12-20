package media

import (
	"context"
	"fmt"
	"path/filepath"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/proto/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/database"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
)

func (s ServiceHandler) GetDirectory(ctx context.Context, req *connect.Request[media.GetDirectoryRequest]) (*connect.Response[media.GetDirectoryResponse], error) {

	res := media.GetDirectoryResponse{
		Path:          req.Msg.Path,
		Name:          filepath.Base(req.Msg.Path),
		ChildrenPaths: make([]string, 0),
		VideoIds:      make([]int32, 0),
	}

	if req.Msg.Path == "/" || req.Msg.Path == "" {
		for _, rootDirectory := range config.Config().MediaDirectories {
			res.ChildrenPaths = append(res.ChildrenPaths, rootDirectory.Path)
		}

		return connect.NewResponse(&res), nil
	}

	dir, _ := filemanager.ReadDirectory(req.Msg.Path)

	for _, childPath := range dir.ChildrenPaths {
		res.ChildrenPaths = append(res.ChildrenPaths, childPath)
	}

	for _, video := range dir.Videos {
		var videoEntry database.Video

		if err := database.Database().Where(&database.Video{
			DirectoryPath: video.DirectoryPath(),
			Filename:      video.Filename(),
		}).
			First(&videoEntry).Error; err != nil {
			return nil, fmt.Errorf("Error getting video entry: %v", err)
		}

		res.VideoIds = append(res.VideoIds, int32(videoEntry.ID))
	}

	return connect.NewResponse(&res), nil
}
