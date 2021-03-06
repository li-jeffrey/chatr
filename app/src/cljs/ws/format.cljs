(ns ws.format
  (:require [cljs.tools.reader :as reader]))


(defprotocol Format
  "Protocol used to define encoding format for socket messages."
  (read  [formatter string])
  (write [formatter value]))

(def json
  "Read and write data encoded in JSON."
  (reify Format
    (read  [_ s] (js->clj (js/JSON.parse s) :keywordize-keys true))
    (write [_ v] (js/JSON.stringify (clj->js v)))))

(def edn
  "Read and write data serialized as EDN."
  (reify Format
    (read [_ s] (reader/read-string s))
    (write [_ v] (prn-str v))))