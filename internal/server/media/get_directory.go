package media

import (
	"context"
	"fmt"
	"path/filepath"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/proto/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/db"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
)

func (s ServiceHandler) GetDirectory(ctx context.Context, req *connect.Request[media.GetDirectoryRequest]) (*connect.Response[media.GetDirectoryResponse], error) {
	if req.Msg.Path == "/" || req.Msg.Path == "" {
		res := media.GetDirectoryResponse{
			Directories: make([]*media.Directory, 0),
			Videos:      make([]*media.Video, 0),
		}

		for _, rootDirectory := range config.Config().Media.RootDirectories {
			res.Directories = append(res.Directories, &media.Directory{
				Path: rootDirectory.Path,
				Name: filepath.Base(rootDirectory.Path),
			})
		}

		return connect.NewResponse(&res), nil
	}

	res := media.GetDirectoryResponse{
		Directories: make([]*media.Directory, 0),
	}

	dir, _, _ := filemanager.ReadDirectory(req.Msg.Path)

	for _, child := range dir.Children() {
		res.Directories = append(res.Directories, &media.Directory{
			Path: child.Path(),
			Name: filepath.Base(child.Path()),
		})
	}

	for _, video := range dir.VideoFiles() {
		var videoEntry db.Video

		if err := db.DB().Where(&db.Video{
			DirectoryPath: video.DirectoryPath(),
			Filename:      video.Filename(),
		}).
			Preload("Subtitles").
			Preload("Subtitles.Segments").
			First(&videoEntry).Error; err != nil {
			return nil, fmt.Errorf("Error getting video entry: %v", err)
		}

		res.Videos = append(res.Videos, &media.Video{
			Id:        int32(videoEntry.ID),
			BaseName:  video.Basename(),
			Extension: video.Extension(),
		})
	}

	return connect.NewResponse(&res), nil
}
