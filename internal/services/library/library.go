package library

import (
	"context"

	"github.com/nandesh-dev/subtle/internal/filemanager"
	"github.com/nandesh-dev/subtle/internal/pb/library"
)

type LibraryServiceServer struct {
	library.UnsafeLibraryServiceServer
}

func (s *LibraryServiceServer) GetMedia(ctx context.Context, req *library.GetMediaRequest) (*library.GetMediaResponse, error) {
	media, _ := filemanager.ReadDirectory("./media")

	var loop func(filemanager.Directory) *library.Directory
	loop = func(dir filemanager.Directory) *library.Directory {
		children := make([]*library.Directory, len(dir.Children()))

		for i, child := range dir.Children() {
			children[i] = loop(child)
		}

		videoFiles, _ := dir.VideoFiles()

		videos := make([]*library.Video, len(videoFiles))

		for i, _ := range videoFiles {
			videos[i] = &library.Video{
				Id: "1",
			}
		}

		return &library.Directory{
			Children: children,
			Videos:   videos,
		}
	}

	return &library.GetMediaResponse{
		Directories: []*library.Directory{loop(*media)},
	}, nil
}
