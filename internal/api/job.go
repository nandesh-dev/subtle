package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/nandesh-dev/subtle/generated/ent/jobschema"
	"github.com/nandesh-dev/subtle/generated/proto/messages"
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
		sequenceNumber = 1
		name = "Extract"
		description = "Extract subtitle from the videos, decode it and store them in database."
		break
	case jobschema.CodeFormat.String():
		sequenceNumber = 1
		name = "Format"
		description = "Format the subtitles in database."
		break
	case jobschema.CodeExport.String():
		sequenceNumber = 1
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
