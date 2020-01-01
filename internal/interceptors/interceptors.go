package interceptors

import (
	"github.com/valyala/fasthttp"
)

var Interceptors = []Interceptor{
	LoggingInterceptor{},
}

// Interceptor is for intercepting requests before and after they are handled by a
// Controller. If at anytime an interceptor returns false, then the request returns immediately.
type Interceptor interface {
	Before(ctx *fasthttp.RequestCtx) bool
	After(ctx *fasthttp.RequestCtx) bool
}
