package main

import (
	"chatr/internal/controllers"
	"chatr/internal/logger"
	"flag"
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var log = logger.GetLogger("main")
var port = flag.Int("port", 8080, "Http port for the application")

func main() {
	flag.Parse()
	router := fasthttprouter.New()
	controllers.RegisterEndpoints(router)

	log.Info("Listening on 127.0.0.1:%v", port)
	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%v", port), router.Handler))
}
