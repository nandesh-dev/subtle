package api

import (
	"context"
	"fmt"
	"path/filepath"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/ent/videoschema"
	"github.com/nandesh-dev/subtle/generated/proto/web"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
)

func (h WebServiceHandler) GetDirectory(ctx context.Context, req *connect.Request[web.GetDirectoryRequest]) (*connect.Response[web.GetDirectoryResponse], error) {
	directory, err := filemanager.ReadDirectory(req.Msg.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory, err: %v", err)
	}

	childrenDirectoryNames := make([]string, 0, len(directory.ChildrenPaths))
	for _, childDirectoryPath := range directory.ChildrenPaths {
		childrenDirectoryNames = append(childrenDirectoryNames, filepath.Base(childDirectoryPath))
	}

	videoFilepaths := make([]string, 0, len(directory.Videos))
	for _, video := range directory.Videos {
		videoFilepaths = append(videoFilepaths, video.Filepath())
	}

	videoEntries, err := h.db.VideoSchema.Query().Where(videoschema.FilepathIn(videoFilepaths...)).All(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get videos from database, err: %v", err)
	}

	videoIds := make([]int32, 0, len(videoEntries))
	for _, videoEntry := range videoEntries {
		videoIds = append(videoIds, int32(videoEntry.ID))
	}

	return connect.NewResponse(&web.GetDirectoryResponse{
		ChildrenDirectoryNames: childrenDirectoryNames,
		VideoIds:               videoIds,
	}), nil
}
