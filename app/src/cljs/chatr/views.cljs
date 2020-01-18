(ns chatr.views
  (:require
    [re-frame.core :as re-frame]
    [re-com.core :as re-com]
    [chatr.subs :as subs]
    [chatr.events :as events]
    [reagent.core :as reagent]))

;; Helper methods
(defn submit-question [_ question]
  (re-frame/dispatch [::events/submit-question question]))

(defn submit-answer [assigned answer]
  (let [id (:ID (last assigned))]
    (re-frame/dispatch [::events/submit-answer id answer])))

(defn submission->alert [{:keys [ID Question Answer]}]
  {
   :id         ID
   :alert-type :none
   :heading    Question
   :body       Answer
   })

;; Components
(defn title-bar [& {:keys [disabled?]
                    :or   {disabled? false}}]
  (let [session-id @(re-frame/subscribe [::subs/session-id])
        disconnect-reason @(re-frame/subscribe [::subs/disconnect-reason])]
    [re-com/h-box
     :width "100%"
     :children [[re-com/gap :size "20px"]
                [re-com/title
                 :label "Chatter"
                 :level :level1
                 :class (if disabled?
                          "main-title alert-state"
                          "main-title")]
                [re-com/md-icon-button
                 :md-icon-name "zmdi-alert-triangle"
                 :tooltip "Disconnected. Click to reconnect"
                 :tooltip-position :below-right
                 :class (if disabled?
                          "alert-state"
                          "hidden")
                 :on-click #(re-frame/dispatch [::events/initialize-ws])
                 :style {:margin-top "1.4em" :margin-left "10px"}]
                [re-com/label
                 :label disconnect-reason
                 :class (if (empty? disconnect-reason)
                          "hidden"
                          "alert-state")
                 :style {:margin-top "2.5em" :margin-left "10px"}]
                [re-com/gap :size "1"]
                [re-com/title
                 :label "SessionID"
                 :level :level4
                 :style {:margin-top "2.5em" :margin-right "10px"}]
                [re-com/input-text
                 :model session-id
                 :on-change #()
                 :change-on-blur? false
                 :width "auto"
                 :style {:margin-top "2em"}]
                [re-com/gap :size "20px"]]]))

(defn chat-panel-heading [& {:keys [label]}]
  [re-com/title
   :label label
   :level :level2
   :class "chat-panel-heading"])

(defn chat-panel-body [& {:keys [items]}]
  [re-com/alert-list
   :alerts (map submission->alert items)
   :on-close #()
   :style {:height "60vh" :border "1px solid #AAAAAA"}])

(defn chat-input [& {:keys [placeholder on-submit disabled?]}]
  (let [text-val (reagent/atom nil)]
    [re-com/h-box
     :width "auto"
     :children [[re-com/input-text
                 :width "90%"
                 :model text-val
                 :change-on-blur? false
                 :on-change #(reset! text-val %)
                 :class "chat-panel-input"
                 :disabled? disabled?
                 :placeholder placeholder
                 :attr {:on-key-up (fn [e]
                                     (if (= 13 (.-keyCode e))
                                       (do (on-submit @text-val)
                                           (reset! text-val nil))))}]
                [re-com/md-icon-button
                 :md-icon-name "zmdi-mail-send"
                 :class "chat-panel-btn"
                 :disabled? disabled?
                 :on-click (fn []
                             (on-submit @text-val)
                             (reset! text-val nil))]]]))

(defn chat-panel [& {:keys [label sub-key placeholder on-submit disabled?]
                     :or   {on-submit #(println %)
                            disabled? false}}]
  (let [items @(re-frame/subscribe [sub-key])]
    [re-com/v-box
     :style (merge (re-com/flex-child-style "1")
                   {:padding "0px 20px 0px 20px"})
     :class "chat-panel"
     :children [[chat-panel-heading
                 :label label]
                [chat-panel-body
                 :items items]
                [chat-input
                 :placeholder placeholder
                 :on-submit (partial on-submit items)
                 :disabled? disabled?]]]))

(defn main-panel []
  (let [connected? @(re-frame/subscribe [::subs/connected?])]
    [re-com/v-box
     :height "100%"
     :width "100%"
     :children [[title-bar
                 :disabled? (not connected?)]
                [re-com/line]
                [re-com/h-split
                 :panel-1 [chat-panel
                           :label "Submitted"
                           :disabled? (not connected?)
                           :sub-key ::subs/submissions
                           :placeholder "Submit a question..."
                           :on-submit submit-question]
                 :panel-2 [chat-panel
                           :label "Assigned"
                           :disabled? (not connected?)
                           :sub-key ::subs/assignments
                           :placeholder "Write a response..."
                           :on-submit submit-answer]
                 :size "80vh"]]]))
