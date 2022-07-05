package grpc

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	runtimev1pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/json"
	"time"
)

type GetTenantId interface {
	GetTenantId() string
}

type GetAggregateId interface {
	GetAggregateId() string
}

type GetAggregateType interface {
	GetAggregateType() string
}

//
// LoadEvents
// @Description:
// @receiver a
// @param ctx
// @param request
// @return *runtimev1pb.LoadEventResponse
// @return error
//
func (a *api) LoadEvents(ctx context.Context, request *runtimev1pb.LoadEventRequest) (resp *runtimev1pb.LoadEventResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	_, err := a.do(func() (any, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		in := &eventstorage.LoadEventRequest{
			TenantId:      request.GetTenantId(),
			AggregateId:   request.GetAggregateId(),
			AggregateType: request.GetAggregateType(),
		}

		out, err := a.eventStorage.LoadEvent(ctx, in)
		if err != nil {
			return nil, err
		}

		resp = &runtimev1pb.LoadEventResponse{
			TenantId:      out.TenantId,
			AggregateId:   out.AggregateId,
			AggregateType: out.AggregateType,
			Snapshot:      nil,
			Events:        nil,
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
		return resp, nil
	})

	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	headers.Values["Date"] = time.Now().String()
	if resp != nil && resp.Headers == nil {
		resp.Headers = headers
		return resp, err
	}
	return &runtimev1pb.LoadEventResponse{Headers: headers}, err
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
func (a *api) SaveSnapshot(ctx context.Context, request *runtimev1pb.SaveSnapshotRequest) (resp *runtimev1pb.SaveSnapshotResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	_, err := a.do(func() (any, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkRequest("SaveSnapshot", request); err != nil {
			return nil, err
		}

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
		resp = &runtimev1pb.SaveSnapshotResponse{}
		return resp, err
	})
	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	headers.Values["Date"] = time.Now().String()
	if resp != nil && resp.Headers == nil {
		resp.Headers = headers
		return resp, err
	}
	return &runtimev1pb.SaveSnapshotResponse{Headers: headers}, err
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
func (a *api) ApplyEvent(ctx context.Context, request *runtimev1pb.ApplyEventRequest) (resp *runtimev1pb.ApplyEventResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	out, err := a.do(func() (any, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkRequest("ApplyEvent", request); err != nil {
			return nil, err
		}

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

		out, err := a.eventStorage.ApplyEvent(ctx, in)
		return out, err
	})

	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	if out != nil {
		res := out.(*eventstorage.ApplyEventsResponse)
		headers = a.newResponseHeaders(res.Headers)
	}
	headers.Values["Date"] = time.Now().String()
	return &runtimev1pb.ApplyEventResponse{Headers: headers}, err
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
func (a *api) CreateEvent(ctx context.Context, request *runtimev1pb.CreateEventRequest) (resp *runtimev1pb.CreateEventResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	out, err := a.do(func() (any, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkRequest("CreateEvent", request); err != nil {
			return nil, err
		}

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
		out, err := a.eventStorage.CreateEvent(ctx, in)
		return out, err
	})

	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	if out != nil {
		res := out.(*eventstorage.CreateEventResponse)
		headers = a.newResponseHeaders(res.Headers)
	}
	headers.Values["Date"] = time.Now().String()
	return &runtimev1pb.CreateEventResponse{Headers: headers}, err
}

func (a *api) newResponseHeaders(out *eventstorage.ResponseHeaders) *runtimev1pb.ResponseHeaders {
	headers := &runtimev1pb.ResponseHeaders{
		Status:  runtimev1pb.ResponseStatus(out.Status),
		Message: out.Message,
		Values:  out.Values,
	}
	return headers
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
func (a *api) DeleteEvent(ctx context.Context, request *runtimev1pb.DeleteEventRequest) (resp *runtimev1pb.DeleteEventResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	out, err := a.do(func() (any, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}

		if err := a.checkRequest("DeleteEvent", request); err != nil {
			return nil, err
		}

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
		return a.eventStorage.DeleteEvent(ctx, in)
	})
	headers := NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)
	if out != nil {
		res, _ := out.(*eventstorage.DeleteEventResponse)
		headers = a.newResponseHeaders(res.Headers)
	}
	headers.Values["Date"] = time.Now().String()
	return &runtimev1pb.DeleteEventResponse{Headers: headers}, err
}

func (a *api) GetRelations(ctx context.Context, request *runtimev1pb.GetRelationsRequest) (resp *runtimev1pb.GetRelationsResponse, respErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				respErr = err
			}
		}
	}()

	_, err := a.do(func() (any, error) {
		if err := a.checkEventStorageComponent(); err != nil {
			return nil, err
		}
		if err := a.checkRequest("GetRelations", request); err != nil {
			return nil, err
		}
		in := &eventstorage.GetRelationsRequest{
			TenantId:      request.TenantId,
			AggregateType: request.AggregateType,
			Filter:        request.Filter,
			Sort:          request.Sort,
			PageNum:       request.PageNum,
			PageSize:      request.PageSize,
		}
		out, err := a.eventStorage.GetRelations(ctx, in)
		if err != nil {
			return nil, err
		}
		var relations []*runtimev1pb.RelationDto
		if out != nil && len(out.Data) > 0 {
			for _, item := range out.Data {
				dto := runtimev1pb.RelationDto{
					Id:          item.Id,
					TenantId:    item.TenantId,
					AggregateId: item.AggregateId,
					IsDeleted:   item.IsDeleted,
					TableName:   item.TableName,
					Items:       item.Items,
				}
				relations = append(relations, &dto)
			}
		}
		resp = &runtimev1pb.GetRelationsResponse{
			TotalRows:  out.TotalRows,
			TotalPages: out.TotalPages,
			Filter:     out.Filter,
			Sort:       out.Sort,
			PageNum:    out.PageNum,
			PageSize:   out.PageSize,
			Data:       relations,
			IsFound:    out.IsFound,
			Error:      out.Error,
		}
		return resp, nil
	})
	if resp == nil {
		resp = &runtimev1pb.GetRelationsResponse{Headers: NewResponseHeaders(runtimev1pb.ResponseStatus_SUCCESS, err, nil)}
	}
	resp.Headers.Values["Date"] = time.Now().String()
	return resp, err
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
		EventTime:    e.EventTime.AsTime(),
		Topic:        e.Topic,
		Metadata:     *metadata,
		Relations:    e.Relations,
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

func (a *api) checkEventStorageComponent() error {
	if a.eventStorage == nil {
		return errors.New("EventStorage component not initialized, please check the configuration file。")
	}
	return nil
}

func (a *api) checkRequest(methodName string, request interface{}) error {
	if i, ok := request.(GetTenantId); ok {
		if len(i.GetTenantId()) == 0 {
			return fmt.Errorf("grpc.%s(request) error: request.TenantId is nil", methodName)
		}
	}

	if i, ok := request.(GetAggregateId); ok {
		if len(i.GetAggregateId()) == 0 {
			return fmt.Errorf("grpc.%s(request) error: request.AggregateId is nil", methodName)
		}
	}

	if i, ok := request.(GetAggregateType); ok {
		if len(i.GetAggregateType()) == 0 {
			return fmt.Errorf("grpc.%s(request) error: request.AggregateType is nil", methodName)
		}
	}
	return nil
}

func (a *api) do(fun func() (any, error)) (any, error) {
	return fun()
}

func NewResponseHeaders(status runtimev1pb.ResponseStatus, err error, values map[string]string) *runtimev1pb.ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	if err != nil {
		return NewResponseHeadersError(err, values)
	}
	resp := &runtimev1pb.ResponseHeaders{
		Status:  status,
		Message: "Success",
		Values:  values,
	}
	return resp
}

func NewResponseHeadersError(err error, values map[string]string) *runtimev1pb.ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	resp := &runtimev1pb.ResponseHeaders{
		Status:  runtimev1pb.ResponseStatus_ERROR,
		Message: err.Error(),
		Values:  values,
	}
	return resp
}

func NewResponseHeadersSuccess(values map[string]string) *runtimev1pb.ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	resp := &runtimev1pb.ResponseHeaders{
		Status:  runtimev1pb.ResponseStatus_SUCCESS,
		Message: "Success",
		Values:  values,
	}
	return resp
}
