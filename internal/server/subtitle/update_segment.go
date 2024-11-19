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
		return nil, fmt.Errorf("Error updating segment: %v", err)
	}

	res := subtitle_proto.UpdateSegmentResponse{}
	return connect.NewResponse(&res), nil
}
