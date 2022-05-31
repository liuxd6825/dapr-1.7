package grpc

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	runtimev1pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

//
// LoadEvents
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.LoadEventResponse
// @return error
//
func (a *api) LoadEvents(ctx context.Context, request *runtimev1pb.LoadEventRequest) (*runtimev1pb.LoadEventResponse, error) {
	in := &eventstorage.LoadEventRequest{
		TenantId:    request.GetTenantId(),
		AggregateId: request.GetAggregateId(),
	}

	out, err := a.eventStorage.LoadEvent(ctx, in)
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

		metadata, err := json.Marshal(out.Snapshot.Metadata)
		if err != nil {
			return nil, err
		}
		snapshot := &runtimev1pb.LoadEventResponse_SnapshotDto{
			AggregateData:  *outAggData,
			SequenceNumber: out.Snapshot.SequenceNumber,
			Metadata:       string(metadata),
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
				EventVersion:   item.EventVersion,
				SequenceNumber: item.SequenceNumber,
			}
			events = append(events, event)
		}
	}
	resp.Events = events

	return &resp, nil
}

//
// SaveSnapshot
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.SaveSnapshotResponse
// @return error
//
func (a *api) SaveSnapshot(ctx context.Context, request *runtimev1pb.SaveSnapshotRequest) (*runtimev1pb.SaveSnapshotResponse, error) {

	aggregateData, err := newMapInterface(request.AggregateData)
	if err != nil {
		return nil, err
	}

	metadata, err := newMapString(request.Metadata)
	if err != nil {
		return nil, err
	}

	in := &eventstorage.SaveSnapshotRequest{
		TenantId:         request.GetTenantId(),
		AggregateId:      request.GetAggregateId(),
		AggregateType:    request.GetAggregateType(),
		AggregateData:    *aggregateData,
		AggregateVersion: request.GetAggregateVersion(),
		SequenceNumber:   request.GetSequenceNumber(),
		Metadata:         *metadata,
	}
	_, err = a.eventStorage.SaveSnapshot(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := runtimev1pb.SaveSnapshotResponse{}
	return &resp, nil
}

//
// ApplyEvent
// @Description: 应用领域事件
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.ApplyEventsResponse
// @return error
//
func (a *api) ApplyEvent(ctx context.Context, request *runtimev1pb.ApplyEventRequest) (*runtimev1pb.ApplyEventResponse, error) {
	events, err := newEvents(request.Events)
	if err != nil {
		return nil, err
	}
	in := &eventstorage.ApplyEventsRequest{
		TenantId:      request.TenantId,
		AggregateId:   request.AggregateId,
		AggregateType: request.AggregateType,
		Events:        events,
	}

	_, err = a.eventStorage.ApplyEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	return &runtimev1pb.ApplyEventResponse{}, nil
}

//
// CreateEvent
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.CreateEventResponse
// @return error
//
func (a *api) CreateEvent(ctx context.Context, request *runtimev1pb.CreateEventRequest) (*runtimev1pb.CreateEventResponse, error) {
	events, err := newEvents(request.Events)
	if err != nil {
		return nil, err
	}
	in := &eventstorage.CreateEventRequest{
		TenantId:      request.TenantId,
		AggregateId:   request.AggregateId,
		AggregateType: request.AggregateType,
		Events:        events,
	}
	_, err = a.eventStorage.CreateEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	return &runtimev1pb.CreateEventResponse{}, nil
}

//
// DeleteEvent
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.DeleteEventResponse
// @return error
//
func (a *api) DeleteEvent(ctx context.Context, request *runtimev1pb.DeleteEventRequest) (*runtimev1pb.DeleteEventResponse, error) {
	event, err := newEvent(request.Event)
	if err != nil {
		return nil, err
	}
	in := &eventstorage.DeleteEventRequest{
		TenantId:      request.TenantId,
		AggregateId:   request.AggregateId,
		AggregateType: request.AggregateType,
		Event:         event,
	}
	_, err = a.eventStorage.DeleteEvent(ctx, in)
	if err != nil {
		return nil, err
	}
	return &runtimev1pb.DeleteEventResponse{}, nil
}

func newEvents(eventDtoList []*runtimev1pb.EventDto) (*[]eventstorage.EventDto, error) {
	var events []eventstorage.EventDto
	if eventDtoList == nil {
		return &events, nil
	}
	for _, e := range eventDtoList {
		event, err := newEvent(e)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}
	return &events, nil
}

func newEvent(e *runtimev1pb.EventDto) (*eventstorage.EventDto, error) {
	eventData, err := newMapInterface(e.EventData)
	if err != nil {
		return nil, err
	}
	metadata, err := newMapString(e.Metadata)
	if err != nil {
		return nil, err
	}
	event := &eventstorage.EventDto{
		EventId:      e.EventId,
		CommandId:    e.CommandId,
		EventData:    *eventData,
		EventType:    e.EventType,
		EventVersion: e.EventVersion,
		PubsubName:   e.PubsubName,
		Topic:        e.Topic,
		Metadata:     *metadata,
	}
	return event, nil
}
func newMapInterface(jsonStr string) (*map[string]interface{}, error) {
	mapData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(jsonStr), &mapData); err != nil {
		return nil, err
	}
	return &mapData, nil
}

func newMapString(jsonStr string) (*map[string]string, error) {
	mapData := map[string]string{}
	if err := json.Unmarshal([]byte(jsonStr), &mapData); err != nil {
		return nil, err
	}
	return &mapData, nil
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
