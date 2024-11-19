package subtitle

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	subtitle_proto "github.com/nandesh-dev/subtle/generated/proto/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/database"
)

func (s ServiceHandler) UpdateSegment(ctx context.Context, req *connect.Request[subtitle_proto.UpdateSegmentRequest]) (*connect.Response[subtitle_proto.UpdateSegmentResponse], error) {
	var segmentEntry database.Segment

	if err := database.Database().Where(&database.Segment{ID: int(req.Msg.Id)}).
		First(&segmentEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting video entry: %v", err)
	}

	segmentEntry.StartTime = req.Msg.Start.AsDuration()
	segmentEntry.EndTime = req.Msg.End.AsDuration()
	segmentEntry.Text = req.Msg.New.Text

	if err := database.Database().Save(&segmentEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting segment: %v", err)
	}

	var subtitleEntry database.Subtitle

	if err := database.Database().Where(database.Subtitle{ID: segmentEntry.SubtitleID}).Find(&subtitleEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting subtitle: %v", err)
	}

	if !subtitleEntry.IsFormated {
		if err := database.Database().Save(subtitleEntry); err != nil {
			return nil, fmt.Errorf("Error saving subtitle: %v", err)
		}
	}

	res := subtitle_proto.UpdateSegmentResponse{}
	return connect.NewResponse(&res), nil
}
