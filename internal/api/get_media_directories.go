package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/proto/web"
)

func (h WebServiceHandler) GetMediaDirectories(ctx context.Context, req *connect.Request[web.GetMediaDirectoriesRequest]) (*connect.Response[web.GetMediaDirectoriesResponse], error) {
	config, err := h.configFile.Read()
	if err != nil {
		return nil, fmt.Errorf("cannot read config file, err: %v", err)
	}

	paths := make([]string, 0, len(config.Job.Scanning))

	for _, scaningGroup := range config.Job.Scanning {
		paths = append(paths, scaningGroup.DirectoryPath)
	}

	return connect.NewResponse(&web.GetMediaDirectoriesResponse{
		Paths: paths,
	}), nil
}
