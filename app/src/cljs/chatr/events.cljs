(ns chatr.events
  (:require
    [cljs.core.async :refer-macros [go]]
    [re-frame.core :as re-frame]
    [chatr.db :as db]
    [ajax.core :refer [POST]]
    [day8.re-frame.tracing :refer-macros [fn-traced defn-traced]]
    [ws.core :as ws])
  (:import (goog.net XhrIo)))

(re-frame/reg-event-fx
  ::initialize-db
  [(re-frame/inject-cofx :load-session-id)]
  (fn-traced [cofx _]
             {:db       (assoc db/default-db :session-id (:session-id cofx))
              :dispatch [::initialize-ws]}))

;; Websocket requests
(defn ws-handler [{:as   msg
                   :keys [ResponseType]}]
  (let [body (dissoc msg :ResponseType)]
    (case ResponseType
      "Connection" (re-frame/dispatch [::on-connection body])
      "Submission" (re-frame/dispatch [::on-submission body])
      "Assignment" (re-frame/dispatch [::on-assignment body])
      "Disconnection" (re-frame/dispatch [::on-disconnection body])
      (println "Unhandled message: " msg))))

(re-frame/reg-event-db
  ::initialize-ws
  (fn [{:as   db
        :keys [session-id]} _]
    (go (let [url (if (empty? session-id)
                    "ws://localhost:8080/ws"
                    (str "ws://localhost:8080/ws?sessionID=" session-id))
              socket (ws/create :url url
                                :on-message ws-handler
                                :on-close #(re-frame/dispatch [::ws-disconnect]))]
          (re-frame/dispatch [::ws-connect socket])))
    db))

;; Http requests
(defn error-handler [{:keys [status status-text]}]
  (.warn js/console (str status " " status-text)))

(re-frame/reg-event-db
  ::submit-question
  (fn-traced [{:as   db
               :keys [session-id]} [_ text]]
             (POST "http://localhost:8080/api/v1/submit/question"
                   {:url-params    {:sessionID session-id}
                    :body          text
                    :error-handler error-handler
                    })
             db))

(re-frame/reg-event-db
  ::submit-answer
  (fn-traced [{:as   db
               :keys [session-id pending-assignments]} [_ id text]]
             (POST "http://localhost:8080/api/v1/submit/answer"
                   {:url-params    {:id id :sessionID session-id}
                    :body          text
                    :error-handler error-handler})
             (assoc db :pending-assignments (dissoc pending-assignments id))))

;; Websocket events
(re-frame/reg-event-db
  ::ws-connect
  (fn-traced [db [_ socket]]
             (assoc db :ws socket)))

(def session-id->local-store (re-frame/after db/save-session-id))

(re-frame/reg-event-db
  ::ws-disconnect
  (fn [db _]
    (assoc db :connected? false)))

(re-frame/reg-event-db
  ::on-connection
  [session-id->local-store]
  (fn-traced [db [_ {:keys [SessionID]}]]
             (assoc db :session-id SessionID
                       :connected? true
                       :disconnect-reason "")))

(re-frame/reg-event-db
  ::on-disconnection
  (fn-traced [db [_ {:keys [Reason]}]]
             (assoc db :disconnect-reason Reason)))

(re-frame/reg-event-db
  ::on-submission
  (fn-traced [{:as db :keys [submissions]} [_ {:as sub :keys [ID]}]]
             (assoc db :submissions (assoc submissions ID sub))))

(re-frame/reg-event-db
  ::on-assignment
  (fn-traced [{:as db :keys [pending-assignments completed-assignments]}
              [_ {:as sub :keys [ID Answer]}]]
             (if (empty? Answer)
               (assoc db :pending-assignments (assoc pending-assignments ID sub))
               (assoc db :completed-assignments (assoc completed-assignments ID sub)))))