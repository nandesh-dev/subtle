package subtitle

import (
	"context"

	"connectrpc.com/connect"
	subtitle_proto "github.com/nandesh-dev/subtle/generated/proto/subtitle"
)

func (s ServiceHandler) GetAvailableExportFormats(ctx context.Context, req *connect.Request[subtitle_proto.GetAvailableExportFormatsRequest]) (*connect.Response[subtitle_proto.GetAvailableExportFormatsResponse], error) {
	res := subtitle_proto.GetAvailableExportFormatsResponse{
		Formats: []string{"srt"},
	}
	return connect.NewResponse(&res), nil
}
