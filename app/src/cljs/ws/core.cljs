(ns ws.core
  (:require [ws.format :as fmt]))

(defn no-op [_])

(defn- read-message [format handler event]
  (->> (.-data event)
       (fmt/read format)
       (handler)))

(defn status [socket]
  "Retrieves the connection status of the socket."
  (condp = (.-readyState socket)
    0 :connecting
    1 :open
    2 :stopping
    3 :stopped))

(defn create
  [& {:keys [url format on-open on-message on-close on-error]
          :or   {format fmt/json
                 on-open  no-op
                 on-close no-op
                 on-error no-op}}]
  (if-let [sock (js/WebSocket. url)]
    (do
      (set! (.-onopen sock) on-open)
      (set! (.-onmessage sock) (partial read-message format on-message))
      (set! (.-onclose sock) on-close)
      (set! (.-onerror sock) on-error)
      sock)
    (throw (js/Error. (str "Web socket connection failed: " url)))))

(defn send
  "Sends data over socket in the specified format."
  ([socket data]
   (send socket data fmt/json))
  ([socket data format]
   (.send socket (fmt/write format data))))

(defn close
  "Closes the socket connection."
  [socket]
  (.close socket))
