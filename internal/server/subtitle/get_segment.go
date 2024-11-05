package subtitle

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	subtitle_proto "github.com/nandesh-dev/subtle/generated/proto/subtitle"
	"github.com/nandesh-dev/subtle/pkgs/db"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (s ServiceHandler) GetSegment(ctx context.Context, req *connect.Request[subtitle_proto.GetSegmentRequest]) (*connect.Response[subtitle_proto.GetSegmentResponse], error) {
	var segmentEntry db.Segment

	if err := db.DB().Where(&db.Segment{ID: int(req.Msg.Id)}).
		First(&segmentEntry).Error; err != nil {
		return nil, fmt.Errorf("Error getting video entry: %v", err)
	}

	res := subtitle_proto.GetSegmentResponse{
		Start:    durationpb.New(segmentEntry.StartTime),
		End:      durationpb.New(segmentEntry.EndTime),
		Original: &subtitle_proto.OriginalSegment{},
		New: &subtitle_proto.NewSegment{
			Text: segmentEntry.Text,
		},
	}

	if segmentEntry.OriginalText == "" {
		res.Original.Image = segmentEntry.OriginalImage
	} else {
		res.Original.Text = &segmentEntry.OriginalText
	}

	return connect.NewResponse(&res), nil
}
