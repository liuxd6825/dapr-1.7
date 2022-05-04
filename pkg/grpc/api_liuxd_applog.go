package grpc

import (
	"context"
	runtimev1pb "github.com/dapr/dapr/pkg/proto/runtime/v1"
)

func (a *api) WriteEventLog(ctx context.Context, request *runtimev1pb.WriteEventLogRequest) (*runtimev1pb.WriteEventLogResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *api) UpdateEventLog(ctx context.Context, request *runtimev1pb.UpdateEventLogRequest) (*runtimev1pb.UpdateEventLogResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *api) GetEventLogByCommandId(ctx context.Context, request *runtimev1pb.GetEventLogByCommandIdRequest) (*runtimev1pb.GetEventLogByCommandIdResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *api) WriteAppLog(ctx context.Context, request *runtimev1pb.WriteAppLogRequest) (*runtimev1pb.WriteAppLogResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *api) UpdateAppLog(ctx context.Context, request *runtimev1pb.UpdateAppLogRequest) (*runtimev1pb.UpdateAppLogResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *api) GetAppLogById(ctx context.Context, request *runtimev1pb.GetAppLogByIdRequest) (*runtimev1pb.GetAppLogByIdResponse, error) {
	//TODO implement me
	panic("implement me")
}
