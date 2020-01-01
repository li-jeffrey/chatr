(ns chatr.runner
    (:require [doo.runner :refer-macros [doo-tests]]
              [chatr.core-test]))

(doo-tests 'chatr.core-test)
