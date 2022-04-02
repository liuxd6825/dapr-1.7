package http

import (
	"encoding/json"
	"github.com/dapr/components-contrib/liuxd/applog"
	"github.com/valyala/fasthttp"
)

func (a *api) constructLoggerEndpoints() []Endpoint {
	return []Endpoint{
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "logger/event-log/create",
			Version: apiVersionV1,
			Handler: a.writeEventLog,
		},
		{
			Methods: []string{fasthttp.MethodPut},
			Route:   "logger/event-log/update",
			Version: apiVersionV1,
			Handler: a.updateEventLog,
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "logger/event-log/{commandId}",
			Version: apiVersionV1,
			Handler: a.getEventLogByCommandId,
		},
		{
			Methods: []string{fasthttp.MethodPost},
			Route:   "logger/app-log/create",
			Version: apiVersionV1,
			Handler: a.writeEventLog,
		},
		{
			Methods: []string{fasthttp.MethodPut},
			Route:   "logger/app-log/update",
			Version: apiVersionV1,
			Handler: a.updateEventLog,
		},
		{
			Methods: []string{fasthttp.MethodGet},
			Route:   "logger/app-log/{id}",
			Version: apiVersionV1,
			Handler: a.getEventLogByCommandId,
		},
	}
}

func (a *api) writeEventLog(ctx *fasthttp.RequestCtx) {
	data := &applog.WriteEventLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.appLogger.WriteEventLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) updateEventLog(ctx *fasthttp.RequestCtx) {
	data := &applog.UpdateEventLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.appLogger.UpdateEventLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) getEventLogByCommandId(ctx *fasthttp.RequestCtx) {
	tenantId := ctx.UserValue("tenantId").(string)
	pubAppId := ctx.UserValue("pubAppId").(string)
	subAppId := ctx.UserValue("subAppId").(string)
	commandId := ctx.UserValue("commandId").(string)
	req := &applog.GetEventLogByCommandIdRequest{
		TenantId:  tenantId,
		PubAppId:  pubAppId,
		SubAppId:  subAppId,
		CommandId: commandId,
	}
	respData, err := a.appLogger.GetEventLogByCommandId(ctx, req)
	setResponseData(ctx, respData, err)
}

func (a *api) writeAppLog(ctx *fasthttp.RequestCtx) {
	data := &applog.WriteEventLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.appLogger.WriteEventLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) updateAppLog(ctx *fasthttp.RequestCtx) {
	data := &applog.UpdateEventLogRequest{}
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		setResponseData(ctx, nil, err)
		return
	}
	respData, err := a.appLogger.UpdateEventLog(ctx, data)
	setResponseData(ctx, respData, err)
}

func (a *api) getAppLogById(ctx *fasthttp.RequestCtx) {
	tenantId := ctx.UserValue("tenantId").(string)
	pubAppId := ctx.UserValue("pubAppId").(string)
	subAppId := ctx.UserValue("subAppId").(string)
	commandId := ctx.UserValue("commandId").(string)
	req := &applog.GetEventLogByCommandIdRequest{
		TenantId:  tenantId,
		PubAppId:  pubAppId,
		SubAppId:  subAppId,
		CommandId: commandId,
	}
	respData, err := a.appLogger.GetEventLogByCommandId(ctx, req)
	setResponseData(ctx, respData, err)
}
