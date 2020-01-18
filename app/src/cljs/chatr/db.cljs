(ns chatr.db
  (:require [re-frame.core :as re-frame]))

(def default-db
  {:session-id nil
   :ws nil
   :connected? false
   :disconnect-reason ""
   :submissions (array-map)
   :pending-assignments (array-map)
   :completed-assignments {}})

(def ls-key "chatr/session-id")

(defn load-session-id []
  (.getItem js/localStorage ls-key))

(defn save-session-id [db]
  (when (empty? (load-session-id))
    (.setItem js/localStorage ls-key (:session-id db))))

(re-frame/reg-cofx
  :load-session-id
  (fn [cofx _]
    (assoc cofx :session-id (load-session-id))))