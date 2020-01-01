package interceptors

import (
	"chatr/internal/logger"

	"github.com/valyala/fasthttp"
)

var log = logger.GetLogger("interceptors")

type LoggingInterceptor struct{}

func (l LoggingInterceptor) Before(ctx *fasthttp.RequestCtx) bool {
	log.Debug("%s %q %s", ctx.Method(), ctx.RequestURI(), ctx.PostBody())
	return true
}

func (l LoggingInterceptor) After(ctx *fasthttp.RequestCtx) bool {
	log.Debug("Status %v %s", ctx.Response.StatusCode(), ctx.Response.Body())
	return true
}
