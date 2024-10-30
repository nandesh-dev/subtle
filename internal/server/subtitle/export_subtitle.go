package subtitle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"connectrpc.com/connect"
	subtitle_proto "github.com/nandesh-dev/subtle/generated/proto/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/db"
	"github.com/nandesh-dev/subtle/pkgs/srt"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
)

func (s ServiceHandler) ExportSubtitle(ctx context.Context, req *connect.Request[subtitle_proto.ExportSubtitleRequest]) (*connect.Response[subtitle_proto.ExportSubtitleResponse], error) {
	var subtitleEntry db.Subtitle
	if err := db.DB().Where(&db.Subtitle{ID: int(req.Msg.SubtitleId)}).Preload("Segments").First(&subtitleEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting subtitle entry from database: %v", err)
	}

	if req.Msg.ExportFormat != "srt" {
		return nil, fmt.Errorf("Invalid export format: %v", req.Msg.ExportFormat)
	}

	sub := subtitle.NewTextSubtitle()

	for _, segmentEntry := range subtitleEntry.Segments {
		sub.AddSegment(*subtitle.NewTextSegment(segmentEntry.StartTime, segmentEntry.EndTime, segmentEntry.Text))
	}

	out := srt.EncodeSubtitle(*sub)

	path := filepath.Join(req.Msg.ExportDirectoryPath, req.Msg.ExportFilename+".srt")
	if err := os.WriteFile(path, []byte(out), 0644); err != nil {
		return nil, fmt.Errorf("Error writting file to disk: %v", err)
	}

	return &connect.Response[subtitle_proto.ExportSubtitleResponse]{}, nil
}
