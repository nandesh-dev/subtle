package subtitle

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	subtitle_proto "github.com/nandesh-dev/subtle/generated/proto/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/db"
)

func (s ServiceHandler) GetSubtitle(ctx context.Context, req *connect.Request[subtitle_proto.GetSubtitleRequest]) (*connect.Response[subtitle_proto.GetSubtitleResponse], error) {
	var subtitleEntry db.Subtitle

	if err := db.DB().Where(&db.Subtitle{ID: int(req.Msg.Id)}).
		Preload("Segments").
		First(&subtitleEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting video entry: %v", err)
	}

	res := subtitle_proto.GetSubtitleResponse{
		Id:           req.Msg.Id,
		Title:        subtitleEntry.Title,
		Language:     subtitleEntry.Language,
		IsProcessing: false,
		Import: &subtitle_proto.Import{
			IsExternal: subtitleEntry.ImportIsExternal,
			Format:     subtitleEntry.ImportFormat,
		},
		SegmentIds: make([]int32, 0),
	}

	if subtitleEntry.ExportPath != "" {
		res.Export = &subtitle_proto.Export{
			Path:   subtitleEntry.ExportPath,
			Format: subtitleEntry.ExportFormat,
		}
	}

	for _, segmentEntry := range subtitleEntry.Segments {
		res.SegmentIds = append(res.SegmentIds, int32(segmentEntry.ID))
	}

	return connect.NewResponse(&res), nil
}
