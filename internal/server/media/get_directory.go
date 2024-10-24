package media

import (
	"context"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/proto/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
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
		})
	}

	for _, video := range dir.VideoFiles() {
		res.Videos = append(res.Videos, &media.Video{
			Name:      video.Basename(),
			Extension: video.Extension(),
		})
	}

	return connect.NewResponse(&res), nil
}

func compileResponseDirectory(directory filemanager.Directory) *media.Directory {
	return &media.Directory{
		Path: directory.Path(),
	}
}
