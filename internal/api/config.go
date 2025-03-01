package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/proto/messages"
)

func (h WebServiceHandler) GetConfig(ctx context.Context, req *connect.Request[messages.GetConfigRequest]) (*connect.Response[messages.GetConfigResponse], error) {
	configString, err := h.configFile.ReadString()
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&messages.GetConfigResponse{
		Config: configString,
	}), nil
}

func (h WebServiceHandler) UpdateConfig(ctx context.Context, req *connect.Request[messages.UpdateConfigRequest]) (*connect.Response[messages.UpdateConfigResponse], error) {
	if err := h.configFile.WriteString(req.Msg.Config); err != nil {
		return nil, err
	}

	return connect.NewResponse(&messages.UpdateConfigResponse{}), nil
}
