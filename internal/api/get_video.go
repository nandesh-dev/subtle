package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/ent/videoschema"
	"github.com/nandesh-dev/subtle/generated/proto/web"
)

func (h WebServiceHandler) GetVideo(ctx context.Context, req *connect.Request[web.GetVideoRequest]) (*connect.Response[web.GetVideoResponse], error) {
	videoEntry, err := h.db.VideoSchema.Query().Where(videoschema.IDEQ(int(req.Msg.Id))).Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get video entry from database, err: %v", err)
	}

	subtitleEntries, err := videoEntry.QuerySubtitles().All(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subtitle entries from database, err: %v", err)
	}

	subtitleIds := make([]int32, 0, len(subtitleEntries))
	for _, subtitleEntry := range subtitleEntries {
		subtitleIds = append(subtitleIds, int32(subtitleEntry.ID))
	}

	return connect.NewResponse(&web.GetVideoResponse{
		Filepath:    videoEntry.Filepath,
		SubtitleIds: subtitleIds,
	}), nil
}
