package api

import (
	"context"
	"fmt"
	"strconv"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/ent/videoschema"
	"github.com/nandesh-dev/subtle/generated/proto/messages"
	"github.com/nandesh-dev/subtle/pkgs/filemanager"
)

func (h WebServiceHandler) GetRootDirectoryPaths(ctx context.Context, req *connect.Request[messages.GetRootDirectoryPathsRequest]) (*connect.Response[messages.GetRootDirectoryPathsResponse], error) {
	config, err := h.configFile.Read()
	if err != nil {
		return nil, fmt.Errorf("cannot read config file, err: %v", err)
	}

	paths := make([]string, 0, len(config.Job.Scanning))

	for _, scaningGroup := range config.Job.Scanning {
		paths = append(paths, scaningGroup.DirectoryPath)
	}

	return connect.NewResponse(&messages.GetRootDirectoryPathsResponse{
		Paths: paths,
	}), nil
}

func (h WebServiceHandler) ReadDirectory(ctx context.Context, req *connect.Request[messages.ReadDirectoryRequest]) (*connect.Response[messages.ReadDirectoryResponse], error) {
	directory, err := filemanager.ReadDirectory(req.Msg.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory, err: %v", err)
	}

	videoFilepaths := make([]string, 0, len(directory.Videos))
	for _, video := range directory.Videos {
		videoFilepaths = append(videoFilepaths, video.Filepath())
	}

	return connect.NewResponse(&messages.ReadDirectoryResponse{
		ChildrenDirectoryPaths: directory.ChildrenPaths,
    VideoPaths: videoFilepaths,
	}), nil
}

func (h WebServiceHandler) SearchVideo(ctx context.Context, req *connect.Request[messages.SearchVideoRequest]) (*connect.Response[messages.SearchVideoResponse], error) {
	videoEntry, err := h.db.VideoSchema.Query().Where(videoschema.Filepath(req.Msg.Path)).Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get video entry from database, err: %v", err)
	}

	subtitleEntries, err := videoEntry.QuerySubtitles().All(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subtitle entries from database, err: %v", err)
	}

	subtitleIds := make([]string, 0, len(subtitleEntries))
	for _, subtitleEntry := range subtitleEntries {
		subtitleIds = append(subtitleIds,  strconv.FormatInt(int64(subtitleEntry.ID), 10))
	}

	return connect.NewResponse(&messages.SearchVideoResponse{
    Id: strconv.FormatInt(int64(videoEntry.ID), 10),
    Path: videoEntry.Filepath,
    SubtitleIds: subtitleIds,
	}), nil
}
