package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/proto/web"
)

func (h WebServiceHandler) UpdateConfig(ctx context.Context, req *connect.Request[web.UpdateConfigRequest]) (*connect.Response[web.UpdateConfigResponse], error) {
	if err := h.configFile.WriteString(req.Msg.UpdatedConfig); err != nil {
		return nil, err
	}

	return connect.NewResponse(&web.UpdateConfigResponse{}), nil
}
