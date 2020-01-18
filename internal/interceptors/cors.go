package interceptors

import "github.com/valyala/fasthttp"

var corsHeader = []byte("Access-Control-Allow-Origin")
var corsHeaderValue = []byte("*")

type CorsInterceptor struct{}

func (c CorsInterceptor) Before(ctx *fasthttp.RequestCtx) bool {
	return true
}

func (c CorsInterceptor) After(ctx *fasthttp.RequestCtx) bool {
	ctx.Response.Header.AddBytesKV(corsHeader, corsHeaderValue)
	return true
}
