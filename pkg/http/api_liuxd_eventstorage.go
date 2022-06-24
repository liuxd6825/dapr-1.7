package http

import (
	"encoding/json"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
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
			Route:   "event-storage/events/{tenantId}/{id}",
			Version: apiVersionV1,
			Handler: a.getEventById,
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "event-storage/events/apply-events",
			Version: apiVersionV1,
			Handler: a.applyEvents,
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "event-storage/events/create-aggregate",
			Version: apiVersionV1,
			Handler: a.createEvent,
		},
		/*		{
				Methods: []string{fasthttp.MethodGet},
				Route:   "event-storage/aggregates/{tenantId}/{id}",
				Version: apiVersionV1,
				Handler: a.getAggregateById,
			},*/
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "event-storage/snapshot/save",
			Version: apiVersionV1,
			Handler: a.saveSnapshot,
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "event-storage/relations/tenants/{tenantId}/types/{aggregateType}",
			Version: apiVersionV1,
			Handler: a.getRelations,
		},
	}
}

/*func (a *api) getAggregateById(ctx *fasthttp.RequestCtx) {
	if !a.check(ctx) {
		return
	}
	tenantId := ctx.UserValue("tenantId").(string)
	id := ctx.UserValue("id").(string)
	req := &eventstorage.ExistAggregateRequest{
		TenantId:    tenantId,
		AggregateId: id,
	}
	respData, err := a.eventStorage.ExistAggregate(ctx, req)
	setResponseData(ctx, respData, err)
}*/

func (a *api) getRelations(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()
	
	tenantId, ok, _ := getUserValue(ctx, "tenantId")
	if !ok {
		setResponseData(ctx, nil, errors.New("/tenants/{tenantId}"))
		return
	}
	aggregateType, ok, _ := getUserValue(ctx, "aggregateType")
	if !ok {
		setResponseData(ctx, nil, errors.New("/types/{aggregateType}"))
		return
	}

	filter, _, _ := getQueryArgsString(ctx, "filter", "")
	sort, _, _ := getQueryArgsString(ctx, "sort", "")
	pageNum, _, _ := getQueryArgsUint(ctx, "pageNum", 0)
	pageSize, _, _ := getQueryArgsUint(ctx, "pageSize", 20)

	query := &eventstorage.GetRelationsRequest{
		TenantId:      tenantId,
		Filter:        filter,
		AggregateType: aggregateType,
		Sort:          sort,
		PageNum:       pageNum,
		PageSize:      pageSize,
	}
	respData, err := a.eventStorage.GetRelations(ctx, query)
	setResponseData(ctx, respData, err)
}

func (a *api) saveSnapshot(ctx *fasthttp.RequestCtx) {
	if !a.check(ctx) {
		return
	}
	data := eventstorage.SaveSnapshotRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.eventStorage.SaveSnapshot(ctx, &data)
	setResponseData(ctx, respData, err)
}

func (a *api) getEventById(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()
	if !a.check(ctx) {
		return
	}
	tenantId := ctx.UserValue("tenantId").(string)
	id := ctx.UserValue("id").(string)
	data := eventstorage.LoadEventRequest{
		TenantId:    tenantId,
		AggregateId: id,
	}

	respData, err := a.eventStorage.LoadEvent(ctx, &data)
	setResponseData(ctx, respData, err)
}

func (a *api) applyEvents(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()
	if !a.check(ctx) {
		return
	}
	data := eventstorage.ApplyEventsRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := a.eventStorage.ApplyEvent(ctx, &data)
	setResponseData(ctx, respData, err)
}

func (a *api) createEvent(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				setResponseData(ctx, nil, err)
			}
		}
	}()

	if !a.check(ctx) {
		return
	}
	data := eventstorage.CreateEventRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}

	respData, err := a.eventStorage.CreateEvent(ctx, &data)
	setResponseData(ctx, respData, err)
}

func (a *api) check(ctx *fasthttp.RequestCtx) bool {
	if a.eventStorage == nil {
		setResponseData(ctx, nil, errors.New("error: api.eventStorage is nil, please check event storage component file. "))
		return false
	}
	return true
}

func setResponseData(ctx *fasthttp.RequestCtx, data interface{}, err error) {
	ctx.SetContentType("application/json")
	if err != nil {
		respErr := &ResponseError{
			Error:         err.Error(),
			AppName:       "dapr",
			ComponentName: "eventstorage",
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

func getUserValue(ctx *fasthttp.RequestCtx, name string) (string, bool, error) {
	var res string
	isFound := false
	value := ctx.UserValue(name)
	if value != nil {
		res = value.(string)
		isFound = true
	}
	return res, isFound, nil
}

func getQueryArgsString(ctx *fasthttp.RequestCtx, name string, defValue string) (string, bool, error) {
	var res string
	queryArgs := ctx.QueryArgs()
	isFound := queryArgs.Has(name)
	if isFound {
		value := queryArgs.Peek(name)
		res = string(value)
		isFound = true
	} else {
		res = defValue
	}
	return res, isFound, nil
}

func getQueryArgsUint(ctx *fasthttp.RequestCtx, name string, defValue uint64) (uint64, bool, error) {
	var res uint64
	s, isFound, err := getQueryArgsString(ctx, name, "")
	if err != nil {
		return 0, false, err
	} else if !isFound {
		res = defValue
	} else {
		res, err = strconv.ParseUint(s, 10, 64)
	}
	return res, isFound, nil
}
