(ns chatr.css
  (:require [garden.def :refer [defstyles]]))

(def blue "#0074D9")
(def aqua "#7FDBFF")
(def navy "#001f3f")
(def red "#ff4136")
(def gray "#aaaaaa")
(def yellow "#ffdc00")

(defstyles screen
           [:.main-title {:color navy}]
           [:.alert-state {:color red}]
           [:.chat-panel-input {:margin "10px 20px 10px 0px" :border (str "1px solid " gray)}]
           [:.chat-panel-btn {:margin "13px 0px 0px 10px"
                              :color blue
                              :hover-color aqua}])
