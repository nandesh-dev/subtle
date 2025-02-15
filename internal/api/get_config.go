package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/proto/web"
)

func (h WebServiceHandler) GetConfig(ctx context.Context, req *connect.Request[web.GetConfigRequest]) (*connect.Response[web.GetConfigResponse], error){
  configString, err := h.configFile.ReadString()
  if err != nil {
    return nil, err
  }

  return connect.NewResponse(&web.GetConfigResponse{
    Config: configString,
  }), nil
}
