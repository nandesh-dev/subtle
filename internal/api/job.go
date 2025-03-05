package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	"entgo.io/ent/dialect/sql"
	"github.com/nandesh-dev/subtle/generated/ent/joblogschema"
	"github.com/nandesh-dev/subtle/generated/ent/jobschema"
	"github.com/nandesh-dev/subtle/generated/proto/messages"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h WebServiceHandler) GetJobs(ctx context.Context, req *connect.Request[messages.GetJobsRequest]) (*connect.Response[messages.GetJobsResponse], error) {
	return connect.NewResponse(&messages.GetJobsResponse{
		JobCodes: []string{
			jobschema.CodeScan.String(),
			jobschema.CodeExtract.String(),
			jobschema.CodeFormat.String(),
			jobschema.CodeExport.String(),
		},
	}), nil
}

func (h WebServiceHandler) GetJob(ctx context.Context, req *connect.Request[messages.GetJobRequest]) (*connect.Response[messages.GetJobResponse], error) {
	var sequenceNumber int32
	var name string
	var description string

	switch req.Msg.Code {
	case jobschema.CodeScan.String():
		sequenceNumber = 1
		name = "Scan"
		description = "Looks for new video files and scan it for subtitle details like language, format, etc."
		break
	case jobschema.CodeExtract.String():
		sequenceNumber = 2
		name = "Extract"
		description = "Extract subtitle from the videos, decode it and store them in database."
		break
	case jobschema.CodeFormat.String():
		sequenceNumber = 3
		name = "Format"
		description = "Format the subtitles in database."
		break
	case jobschema.CodeExport.String():
		sequenceNumber = 4
		name = "Export"
		description = "Export the subtitle stored in database to file."
		break
	default:
		return nil, fmt.Errorf("invalid job code, code: %v", req.Msg.Code)
	}

	jobEntry, err := h.db.JobSchema.Query().
		Where(jobschema.CodeEQ(jobschema.Code(req.Msg.Code))).Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get job from database, err: %w", err)
	}

	return connect.NewResponse(&messages.GetJobResponse{
		Code:           req.Msg.Code,
		SequenceNumber: sequenceNumber,
		Name:           name,
		Description:    description,
		IsRunning:      jobEntry.IsRunning,
		LastRun:        timestamppb.New(jobEntry.LastRun),
	}), nil
}

func (h WebServiceHandler) GetJobLogs(ctx context.Context, req *connect.Request[messages.GetJobLogsRequest]) (*connect.Response[messages.GetJobLogsResponse], error) {
	jobLogQueryBuilder := h.db.JobLogSchema.Query()

	if req.Msg.NewerThanLogId != nil && strings.TrimSpace(*req.Msg.NewerThanLogId) != "" {
		newerThanLogId, err := strconv.Atoi(*req.Msg.NewerThanLogId)
		if err != nil {
			return nil, fmt.Errorf("invalid newer than log id, err: %w", err)
		}

		jobLogQueryBuilder = jobLogQueryBuilder.Where(joblogschema.IDGT(newerThanLogId))
	}

	if req.Msg.OlderThanLogId != nil && strings.TrimSpace(*req.Msg.OlderThanLogId) != "" {
		olderThanLogId, err := strconv.Atoi(*req.Msg.OlderThanLogId)
		if err != nil {
			return nil, fmt.Errorf("invalid older than log id, err: %w", err)
		}

		jobLogQueryBuilder = jobLogQueryBuilder.Where(joblogschema.IDLT(olderThanLogId))
	}

	if req.Msg.Limit != nil {
		jobLogQueryBuilder = jobLogQueryBuilder.Limit(int(*req.Msg.Limit))
	}

	jobLogEntries, err := jobLogQueryBuilder.
		Order(joblogschema.ByID(sql.OrderDesc())).All(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get job log entries from database, err: %w", err)
	}

	ids := make([]string, len(jobLogEntries))
	for i, jobLogEntry := range jobLogEntries {
		ids[i] = strconv.Itoa(jobLogEntry.ID)
	}

	return connect.NewResponse(&messages.GetJobLogsResponse{Ids: ids}), nil
}

func (h WebServiceHandler) GetJobLog(ctx context.Context, req *connect.Request[messages.GetJobLogRequest]) (*connect.Response[messages.GetJobLogResponse], error) {
	id, err := strconv.Atoi(req.Msg.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid job log id, err: %w", err)
	}

	jobLogEntry, err := h.db.JobLogSchema.Query().
		Where(joblogschema.IDEQ(id)).Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get job log from database, err: %w", err)
	}

	jobEntry, err := jobLogEntry.QueryJob().Only(h.ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get job from database, err: %w", err)
	}

	var name string
	switch jobEntry.Code {
	case jobschema.CodeScan:
		name = "Scan"
	case jobschema.CodeExtract:
		name = "Extract"
	case jobschema.CodeFormat:
		name = "Format"
	case jobschema.CodeExport:
		name = "Export"
	default:
		name = "Unknown"
	}

	return connect.NewResponse(&messages.GetJobLogResponse{
		Id:             req.Msg.Id,
		JobCode:        jobEntry.Code.String(),
		JobName:        name,
		StartTimestamp: timestamppb.New(jobLogEntry.StartTimestamp),
		Duration:       durationpb.New(time.Duration(jobLogEntry.Duration)),
	}), nil
}
