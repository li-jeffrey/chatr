(ns chatr.config)

(def debug?
  ^boolean goog.DEBUG)

(def http-base
  (if debug?
    "http://localhost:8080/api/v1"
    "/api/v1"))

(def ws-url
  (if debug?
    "ws://localhost:8080/api/v1/ws"
    (str "wss://" (.-host (.-location js/window)) "/api/v1/ws")))