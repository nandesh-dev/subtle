package api

import (
	"context"
	"fmt"
	"strconv"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/ent/subtitlecuecontentsegmentschema"
	"github.com/nandesh-dev/subtle/generated/ent/subtitlecueoriginalimageschema"
	"github.com/nandesh-dev/subtle/generated/ent/subtitlecueschema"
	"github.com/nandesh-dev/subtle/generated/ent/subtitleschema"
	"github.com/nandesh-dev/subtle/generated/ent/videoschema"
	"github.com/nandesh-dev/subtle/generated/proto/messages"
	"github.com/nandesh-dev/subtle/pkgs/subtitle"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (h WebServiceHandler) CalculateSubtitleStatistics(ctx context.Context, req *connect.Request[messages.CalculateSubtitleStatisticsRequest]) (*connect.Response[messages.CalculateSubtitleStatisticsResponse], error) {
	totalVideoCount, err := h.db.VideoSchema.Query().Count(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get total video count from database, err: %w", err)
	}

	videoWithExportedSubtitleCount, err := h.db.VideoSchema.Query().
		Where(videoschema.HasSubtitlesWith(
			subtitleschema.StageEQ(subtitleschema.StageExported),
		)).Count(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot count video with exported subtitle from database, err: %w", err)
	}

	videoWithFormatedSubtitleCount, err := h.db.VideoSchema.Query().
		Where(videoschema.HasSubtitlesWith(
			subtitleschema.Or(
				subtitleschema.StageEQ(subtitleschema.StageFormated),
				subtitleschema.StageEQ(subtitleschema.StageExported),
			),
		)).Count(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot count video with formated subtitle from database, err: %w", err)
	}

	videoWithExtractedSubtitleCount, err := h.db.VideoSchema.Query().
		Where(videoschema.HasSubtitlesWith(
			subtitleschema.Or(
				subtitleschema.StageEQ(subtitleschema.StageExtracted),
				subtitleschema.StageEQ(subtitleschema.StageFormated),
				subtitleschema.StageEQ(subtitleschema.StageExported),
			),
		)).Count(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot count video with extracted subtitle from database, err: %w", err)
	}

	return connect.NewResponse(&messages.CalculateSubtitleStatisticsResponse{
		TotalVideoCount:                 int32(totalVideoCount),
		VideoWithExtractedSubtitleCount: int32(videoWithExtractedSubtitleCount),
		VideoWithFormatedSubtitleCount:  int32(videoWithFormatedSubtitleCount),
		VideoWithExportedSubtitleCount:  int32(videoWithExportedSubtitleCount),
	}), nil
}

func (h WebServiceHandler) GetSubtitle(ctx context.Context, req *connect.Request[messages.GetSubtitleRequest]) (*connect.Response[messages.GetSubtitleResponse], error) {
	id, err := strconv.Atoi(req.Msg.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid subtitle id, err: %w", err)
	}

	subtitleEntry, err := h.db.SubtitleSchema.Query().
		Where(subtitleschema.IDEQ(id)).Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get video entry from database, err: %w", err)
	}

	subtitleCueEntries, err := subtitleEntry.QueryCues().
		Order(subtitlecueschema.ByTimestampStart()).
		All(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subtitle cue entries from database, err: %w", err)
	}

	subtitleCueIds := make([]string, len(subtitleCueEntries))
	for i, subtitleCueEntry := range subtitleCueEntries {
		subtitleCueIds[i] = strconv.Itoa(subtitleCueEntry.ID)
	}

	stage := messages.SubtitleStage_SUBTITLE_STAGE_UNSPECIFIED
	switch subtitleEntry.Stage {
	case subtitleschema.StageDetected:
		stage = messages.SubtitleStage_SUBTITLE_STAGE_DETECTED
	case subtitleschema.StageExtracted:
		stage = messages.SubtitleStage_SUBTITLE_STAGE_EXTRACTED
	case subtitleschema.StageFormated:
		stage = messages.SubtitleStage_SUBTITLE_STAGE_FORMATED
	case subtitleschema.StageExported:
		stage = messages.SubtitleStage_SUBTITLE_STAGE_EXPORTED
	}

	originalFormat := messages.SubtitleOriginalFormat_SUBTITLE_ORIGINAL_FORMAT_UNSPECIFIED
	switch subtitleEntry.ImportFormat {
	case subtitle.SRT:
		originalFormat = messages.SubtitleOriginalFormat_SUBTITLE_ORIGINAL_FORMAT_SRT
	case subtitle.ASS:
		originalFormat = messages.SubtitleOriginalFormat_SUBTITLE_ORIGINAL_FORMAT_ASS
	case subtitle.PGS:
		originalFormat = messages.SubtitleOriginalFormat_SUBTITLE_ORIGINAL_FORMAT_PGS
	}

	return connect.NewResponse(&messages.GetSubtitleResponse{
		Id:               req.Msg.Id,
		Title:            subtitleEntry.Title,
		Language:         subtitleEntry.Language.String(),
		Stage:            stage,
		IsProcessing:     subtitleEntry.IsProcessing,
		ImportIsExternal: subtitleEntry.ImportIsExternal,
		OriginalFormat:   originalFormat,
		CueIds:           subtitleCueIds,
	}), nil
}

func (h WebServiceHandler) GetSubtitleCue(ctx context.Context, req *connect.Request[messages.GetSubtitleCueRequest]) (*connect.Response[messages.GetSubtitleCueResponse], error) {
	id, err := strconv.Atoi(req.Msg.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid subtitle cue id, err: %w", err)
	}

	subtitleCueEntry, err := h.db.SubtitleCueSchema.Query().
		Where(subtitlecueschema.IDEQ(id)).Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subtitle cue entry from database, err: %w", err)
	}

	segmentEntries, err := subtitleCueEntry.QueryContentSegments().
		All(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subtitle cue segment entries from database, err: %w", err)
	}

	segmentIds := make([]string, len(segmentEntries))
	for i, segmentEntry := range segmentEntries {
		segmentIds[i] = strconv.Itoa(segmentEntry.ID)
	}

	return connect.NewResponse(&messages.GetSubtitleCueResponse{
		Id:         req.Msg.Id,
		Start:      durationpb.New(subtitleCueEntry.TimestampStart),
		End:        durationpb.New(subtitleCueEntry.TimestampEnd),
		SegmentIds: segmentIds,
	}), nil
}

func (h WebServiceHandler) GetSubtitleCueOriginalData(ctx context.Context, req *connect.Request[messages.GetSubtitleCueOriginalDataRequest]) (*connect.Response[messages.GetSubtitleCueOriginalDataResponse], error) {
	subtitleCueId, err := strconv.Atoi(req.Msg.SubtitleCueId)
	if err != nil {
		return nil, fmt.Errorf("invalid subtitle cue id, err: %w", err)
	}

	originalImageEntries, err := h.db.SubtitleCueOriginalImageSchema.Query().
		Where(subtitlecueoriginalimageschema.HasCueWith(subtitlecueschema.IDEQ(subtitleCueId))).
		Order(subtitlecueoriginalimageschema.ByPosition()).
		All(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subtitle original images entries from database, err: %w", err)
	}

	originalImagesData := make([][]byte, len(originalImageEntries))
	for i, originalImageEntry := range originalImageEntries {
		originalImagesData[i] = originalImageEntry.Data
	}

	return connect.NewResponse(&messages.GetSubtitleCueOriginalDataResponse{
		SubtitleCueId: req.Msg.SubtitleCueId,
		//TODO Add original text
    PngEncodedImagesData: originalImagesData,
	}), nil
}

func (h WebServiceHandler) GetSubtitleCueSegment(ctx context.Context, req *connect.Request[messages.GetSubtitleCueSegmentRequest]) (*connect.Response[messages.GetSubtitleCueSegmentResponse], error) {
	id, err := strconv.Atoi(req.Msg.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid subtitle cue segment id, err: %w", err)
	}

	subtitleCueSegmentEntry, err := h.db.SubtitleCueContentSegmentSchema.Query().
		Where(subtitlecuecontentsegmentschema.IDEQ(id)).Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subtitle cue entry from database, err: %w", err)
	}

	return connect.NewResponse(&messages.GetSubtitleCueSegmentResponse{
		Id:   req.Msg.Id,
		Text: subtitleCueSegmentEntry.Text,
	}), nil
}
