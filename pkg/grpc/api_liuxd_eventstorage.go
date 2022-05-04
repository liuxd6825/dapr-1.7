package grpc

import (
	"context"
	"github.com/dapr/components-contrib/liuxd/eventstorage"
	runtimev1pb "github.com/dapr/dapr/pkg/proto/runtime/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

// LoadEvent 是否存在聚合根
func (a *api) LoadEvents(ctx context.Context, request *runtimev1pb.LoadEventRequest) (*runtimev1pb.LoadEventResponse, error) {
	in := &eventstorage.LoadEventRequest{
		TenantId:    request.GetTenantId(),
		AggregateId: request.GetAggregateId(),
	}

	out, err := a.eventStorage.LoadEvents(ctx, in)
	if err != nil {
		return nil, err
	}

	resp := runtimev1pb.LoadEventResponse{
		TenantId:    out.TenantId,
		AggregateId: out.AggregateId,
		Snapshot:    nil,
		Events:      nil,
	}

	if out.Snapshot != nil {
		outAggData, err := mapAsStr(out.Snapshot.AggregateData)
		if err != nil {
			return nil, err
		}

		snapshot := &runtimev1pb.LoadEventResponse_SnapshotDto{
			AggregateData:  *outAggData,
			SequenceNumber: resp.Snapshot.SequenceNumber,
			Metadata:       resp.Snapshot.Metadata,
		}
		resp.Snapshot = snapshot
	}

	events := make([]*runtimev1pb.LoadEventResponse_EventDto, 0)
	if out.Events != nil {
		for _, item := range *out.Events {
			eventData, err := mapAsStr(item.EventData)
			if err != nil {
				return nil, err
			}
			event := &runtimev1pb.LoadEventResponse_EventDto{
				EventId:        item.EventId,
				EventType:      item.EventType,
				EventData:      *eventData,
				EventRevision:  item.EventRevision,
				SequenceNumber: item.SequenceNumber,
			}
			events = append(events, event)
		}
	}
	resp.Events = events

	return &resp, nil
}

func (a *api) SaveSnapshot(ctx context.Context, request *runtimev1pb.SaveSnapshotRequest) (*runtimev1pb.SaveSnapshotResponse, error) {

	aggregateData, err := newMap(request.AggregateData)
	if err != nil {
		return nil, err
	}

	in := &eventstorage.SaveSnapshotRequest{
		TenantId:          request.GetTenantId(),
		AggregateId:       request.GetAggregateId(),
		AggregateType:     request.GetAggregateType(),
		AggregateData:     aggregateData,
		AggregateRevision: request.GetAggregateRevision(),
		SequenceNumber:    request.GetSequenceNumber(),
		Metadata:          request.GetMetadata(),
	}
	_, err = a.eventStorage.SaveSnapshot(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := runtimev1pb.SaveSnapshotResponse{}
	return &resp, nil
}

func (a *api) ExistAggregate(ctx context.Context, request *runtimev1pb.ExistAggregateRequest) (*runtimev1pb.ExistAggregateResponse, error) {
	in := &eventstorage.ExistAggregateRequest{
		TenantId:    request.TenantId,
		AggregateId: request.AggregateId,
	}
	out, err := a.eventStorage.ExistAggregate(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := runtimev1pb.ExistAggregateResponse{
		IsExist: out.IsExist,
	}
	return &resp, nil
}

func (a *api) ApplyEvent(ctx context.Context, request *runtimev1pb.ApplyEventRequest) (*runtimev1pb.ApplyEventResponse, error) {
	eventData, err := newMap(request.EventData)
	if err != nil {
		return nil, err
	}

	in := &eventstorage.ApplyEventRequest{
		TenantId:      request.TenantId,
		Metadata:      request.Metadata,
		CommandId:     request.CommandId,
		EventType:     request.EventType,
		EventId:       request.EventId,
		EventRevision: request.EventRevision,
		EventData:     eventData,
		AggregateId:   request.AggregateId,
		AggregateType: request.AggregateType,
		PubsubName:    request.PubsubName,
		Topic:         request.Topic,
	}
	_, err = a.eventStorage.ApplyEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &runtimev1pb.ApplyEventResponse{}
	return resp, nil
}

func newMap(jsonStr string) (map[string]interface{}, error) {
	mapData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonStr), mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}

func mapAsStr(data map[string]interface{}) (*string, error) {
	if data == nil {
		var empty = ""
		return &empty, nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var res = string(bytes)
	return &res, nil
}

func (a *api) mustEmbedUnimplementedDaprServer() {

}
