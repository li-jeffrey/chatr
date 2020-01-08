package main

import (
	"chatr/internal/controllers"
	"chatr/internal/logger"
	"chatr/internal/store"
	"flag"
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var log = logger.GetLogger("main")
var port = flag.Int("port", 8080, "Http port for the application. Default is 8080")
var storeLocation = flag.String("store.location", "store.bin", "Location of store file. Default is store.bin")

func main() {
	flag.Parse()
	store.EnablePersistence(*storeLocation)
	router := fasthttprouter.New()
	controllers.RegisterEndpoints(router)

	log.Info("Listening on 127.0.0.1:%v", *port)
	log.Fatal("Shutting down: %s", fasthttp.ListenAndServe(fmt.Sprintf(":%v", *port), router.Handler))
}
