package controllers

import (
	"chatr/internal/interceptors"
	"chatr/internal/logger"
	"path"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const prefix string = "/api/v1"

var (
	log                = logger.GetLogger("controllers")
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

// RegisterEndpoints is a function which registers endpoints with the router
func RegisterEndpoints(router *fasthttprouter.Router) {
	handle(router, "POST", "/submit/question", SubmitQuestion)
	handle(router, "POST", "/submit/answer", SubmitAnswer)
	handleWs(router, wsPath)
	handleFS(router, "/", "public")
	handleFS(router, "/js/*filepath", "public/js")
	handleFS(router, "/vendor/*filepath", "public/vendor")
	handleFS(router, "/css/*filepath", "public/css")
	handleNotFound(router)
}

// RequestHandler is a functional interface which consumes a requestCtx and returns a result.
type RequestHandler func(ctx *fasthttp.RequestCtx)

func handle(router *fasthttprouter.Router, method string, requestPath string, handler fasthttp.RequestHandler) {
	fullPath := path.Join(prefix, requestPath)
	router.Handle(method, fullPath, intercept(handler))
	log.Info("Registered handler for %s on %s", method, fullPath)
}

func handleWs(router *fasthttprouter.Router, wsPath string) {
	router.Handle("GET", wsPath, upgradeConnection)
	log.Info("Registered ws connection on %s", wsPath)
}

func handleFS(router *fasthttprouter.Router, path string, relativePath string) {
	router.Handle("GET", path, fasthttp.FSHandler(relativePath, 1))
	log.Info("Registered FSHandler on %s", path)
}

func handleNotFound(router *fasthttprouter.Router) {
	router.NotFound = intercept(func(ctx *fasthttp.RequestCtx) {
		ctx.Error("Not found", fasthttp.StatusNotFound)
	})
}

func intercept(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		for _, interceptor := range interceptors.Interceptors {
			if !interceptor.Before(ctx) {
				return
			}
		}

		handler(ctx)

		for _, interceptor := range interceptors.Interceptors {
			if !interceptor.After(ctx) {
				return
			}
		}
	}
}
