package http

import (
	"encoding/json"
	"net/http"

	eventsourcing "github.com/dapr/components-contrib/eventsourcing/v1"
	"github.com/valyala/fasthttp"
)

type ResponseError struct {
	Error         string `json:"error"`
	AppName       string `json:"appName"`
	ComponentName string `json:"componentName"`
}

func (a *api) constructEventSourcingEndpoints() []Endpoint {
	return []Endpoint{
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "event-sourcing/events/{tenantId}/{id}",
			Version: apiVersionV1,
			Handler: a.getEventById,
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "event-sourcing/events/apply",
			Version: apiVersionV1,
			Handler: a.applyEvent,
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "event-sourcing/aggregates/{tenantId}/{id}",
			Version: apiVersionV1,
			Handler: a.getAggregateById,
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "event-sourcing/snapshot/save",
			Version: apiVersionV1,
			Handler: a.saveSnapshot,
		},
	}
}

func (a *api) getAggregateById(ctx *fasthttp.RequestCtx) {
	tenantId := ctx.UserValue("tenantId").(string)
	id := ctx.UserValue("id").(string)
	req := &eventsourcing.ExistAggregateRequest{
		TenantId:    tenantId,
		AggregateId: id,
	}
	respData, err := a.eventSourcing.ExistAggregate(ctx, req)
	setResponseData(ctx, respData, err)
}

func (a *api) saveSnapshot(ctx *fasthttp.RequestCtx) {
	data := eventsourcing.SaveSnapshotRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.eventSourcing.SaveSnapshot(ctx, &data)
	setResponseData(ctx, respData, err)
}

func (a *api) getEventById(ctx *fasthttp.RequestCtx) {
	tenantId := ctx.UserValue("tenantId").(string)
	id := ctx.UserValue("id").(string)
	data := eventsourcing.LoadEventRequest{
		TenantId:    tenantId,
		AggregateId: id,
	}

	respData, err := a.eventSourcing.LoadEvents(ctx, &data)
	setResponseData(ctx, respData, err)
}

func (a *api) applyEvent(ctx *fasthttp.RequestCtx) {
	data := eventsourcing.ApplyEventRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := a.eventSourcing.ApplyEvent(ctx, &data)
	setResponseData(ctx, respData, err)
}

func setResponseData(ctx *fasthttp.RequestCtx, data interface{}, err error) {
	ctx.SetContentType("application/json")
	if err != nil {
		respErr := &ResponseError{
			Error:         err.Error(),
			AppName:       "dapr",
			ComponentName: "event-sourcing",
		}
		_, _ = ctx.Write(getJsonBytes(respErr))
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.Success("application/json", getJsonBytes(data))
}

func getJsonBytes(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}
