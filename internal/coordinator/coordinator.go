package coordinator

import (
	"bytes"
	"chatr/internal/store"
	"fmt"

	"github.com/valyala/fastjson"
)

var (
	SubscriptionRequest []byte = []byte("Subscription")
	Submission          []byte = []byte("Submission")
	Assignment          []byte = []byte("Assignment")
	BadMessage          []byte = []byte("BAD_MESSAGE")
	InvalidRequest      []byte = []byte("INVALID_REQUEST")
)

func init() {
	store.SubscribeCreations(onSubmissionCreated)
}

func onSubmissionCreated(sub *store.Submission) {
	s := <-selectNextSession()
	s.sendResponse(Assignment, sub)
}

func handleMessage(s *Session, v *fastjson.Value) {
	requestType := v.GetStringBytes("RequestType")
	if requestType == nil {
		s.sendError(InvalidRequest, "Missing request type")
		return
	}

	if bytes.Equal(requestType, SubscriptionRequest) {
		handleSubscriptionRequest(s, v)
	} else {
		s.sendError(InvalidRequest, fmt.Sprintf("Unknown request type %s", requestType))
	}
}

func handleSubscriptionRequest(s *Session, v *fastjson.Value) {
	for _, qid := range v.GetArray("IDs") {
		store.SubscribeSubmission(string(qid.GetStringBytes()), func(sub *store.Submission) {
			s.sendResponse(Submission, sub)
		})
	}
}
