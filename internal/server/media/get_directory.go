package media

import (
	"context"

	"github.com/nandesh-dev/subtle/generated/api/media"
	"github.com/nandesh-dev/subtle/pkgs/config"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
)

func (s *MediaServiceServer) GetDirectory(ctx context.Context, req *media.GetDirectoryRequest) (*media.GetDirectoryResponse, error) {
	if req.Path == "/" || req.Path == "" {
		res := media.GetDirectoryResponse{
			Directories: make([]*media.Directory, 0),
			Videos:      make([]*media.Video, 0),
		}

		for _, rootDirectory := range config.Config().Media.RootDirectories {
			res.Directories = append(res.Directories, &media.Directory{
				Path: rootDirectory.Path,
			})
		}

		return &res, nil
	}

	res := media.GetDirectoryResponse{
		Directories: make([]*media.Directory, 0),
	}

	dir, _, _ := filemanager.ReadDirectory(req.Path)

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

	return &res, nil
}

func compileResponseDirectory(directory filemanager.Directory) *media.Directory {
	return &media.Directory{
		Path: directory.Path(),
	}
}
