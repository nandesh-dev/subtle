package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/ent/subtitleschema"
	"github.com/nandesh-dev/subtle/generated/proto/web"
)

func (h WebServiceHandler) GetGlobalStatistics(ctx context.Context, req *connect.Request[web.GetGlobalStatisticsRequest]) (*connect.Response[web.GetGlobalStatisticsResponse], error){
  totalSubtitleCount, err := h.db.SubtitleSchema.Query().Count(h.ctx)
  if err != nil {
    return nil, fmt.Errorf("Error: %v", err)
  }
  extractedSubtitleCount, err := h.db.SubtitleSchema.Query().Where(subtitleschema.StageEQ(subtitleschema.StageExtracted)).Count(h.ctx)
  if err != nil {
    return nil, fmt.Errorf("Error: %v", err)
  }

  formatedSubtitleCount, err := h.db.SubtitleSchema.Query().Where(subtitleschema.StageEQ(subtitleschema.StageFormated)).Count(h.ctx)
  if err != nil {
    return nil, fmt.Errorf("Error: %v", err)
  }

  exportedSubtitleCount, err := h.db.SubtitleSchema.Query().Where(subtitleschema.StageEQ(subtitleschema.StageExported)).Count(h.ctx)
  if err != nil {
    return nil, fmt.Errorf("Error: %v", err)
  }

  return connect.NewResponse(&web.GetGlobalStatisticsResponse{
    Extracted: int32(extractedSubtitleCount + formatedSubtitleCount + exportedSubtitleCount),
    Formated: int32(formatedSubtitleCount + exportedSubtitleCount),
    Exported: int32(exportedSubtitleCount),
    Total: int32(totalSubtitleCount),
  }), nil
}
