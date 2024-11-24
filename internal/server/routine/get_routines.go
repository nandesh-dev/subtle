package routine

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	routine_proto "github.com/nandesh-dev/subtle/generated/proto/routine"
	"github.com/nandesh-dev/subtle/pkgs/database"
)

func (s ServiceHandler) GetRoutines(ctx context.Context, req *connect.Request[routine_proto.GetRoutinesRequest]) (*connect.Response[routine_proto.GetRoutinesResponse], error) {
	var routineEntries []database.Routine

	if err := database.Database().Find(&routineEntries).Error; err != nil {
		return nil, fmt.Errorf("Error getting routine entries: %v", err)
	}

	res := routine_proto.GetRoutinesResponse{
		Routines: make([]*routine_proto.Routine, 0, len(routineEntries)),
	}

	for _, routineEntry := range routineEntries {
		res.Routines = append(res.Routines, &routine_proto.Routine{
			Name:        routineEntry.Name,
			Description: routineEntry.Description,
			IsRunning:   routineEntry.IsRunning,
		})
	}

	return connect.NewResponse(&res), nil
}
