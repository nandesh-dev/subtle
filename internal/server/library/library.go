package library

import (
	"context"

	"github.com/nandesh-dev/subtle/generated/api/library"
)

type LibraryServiceServer struct {
	library.UnsafeLibraryServiceServer
}

func (s *LibraryServiceServer) GetMedia(ctx context.Context, req *library.GetMediaRequest) (*library.GetMediaResponse, error) {
	return nil, nil
}
