(ns chatr.subs
  (:require
    [re-frame.core :as re-frame]))

(re-frame/reg-sub
  ::submissions
  (fn [db _]
    (vals (:submissions db))))

(re-frame/reg-sub
  ::assignments
  (fn [db _]
    (concat (vals (:completed-assignments db))
            (reverse (vals (:pending-assignments db))))))

(re-frame/reg-sub
  ::connected?
  (fn [db _]
    (:connected? db)))

(re-frame/reg-sub
  ::session-id
  (fn [db _]
    (:session-id db)))

(re-frame/reg-sub
  ::disconnect-reason
  (fn [db _]
    (:disconnect-reason db)))