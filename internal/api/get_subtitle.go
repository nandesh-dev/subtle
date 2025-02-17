package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/ent/subtitleschema"
	"github.com/nandesh-dev/subtle/generated/proto/web"
)

func (h WebServiceHandler) GetSubtitle(ctx context.Context, req *connect.Request[web.GetSubtitleRequest]) (*connect.Response[web.GetSubtitleResponse], error) {
	subtitleEntry, err := h.db.SubtitleSchema.Query().Where(subtitleschema.IDEQ(int(req.Msg.Id))).Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get video entry from database, err: %v", err)
	}

  stage := web.SubtitleStage_DETECTED
  switch subtitleEntry.Stage {
  case subtitleschema.StageExtracted:
    stage = web.SubtitleStage_EXTRACTED
  case subtitleschema.StageFormated:
    stage = web.SubtitleStage_FORMATED
  case subtitleschema.StageExported:
    stage = web.SubtitleStage_EXPORTED
  }

	return connect.NewResponse(&web.GetSubtitleResponse{
    Title: subtitleEntry.Title,
    Language: subtitleEntry.Language.String(),
    Stage: stage,
    IsProcessing: subtitleEntry.IsProcessing,
    ImportIsExternal: subtitleEntry.ImportIsExternal,
    ImportFormat: subtitleEntry.ImportFormat.ToProto(),
	}), nil
}
