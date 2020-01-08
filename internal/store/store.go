package store

import (
	"chatr/internal/logger"
	"fmt"
	"sync"
	"time"

	"github.com/rs/xid"
)

// Submission is a record containing a question-answer pair and its ID.
type Submission struct {
	ID         xid.ID
	Question   []byte
	Answer     []byte
	LastUpdate int64
}

type Callback func(*Submission)

var log = logger.GetLogger("store")

type event interface{}

// events
type submissionCreated struct{}

func (e submissionCreated) String() string {
	return "{SubmissionCreated}"
}

type submissionChanged struct {
	id xid.ID
}

func (s submissionChanged) String() string {
	return fmt.Sprintf("{SubmissionChanged %s}", s.id)
}

// stored state
var mutex sync.Mutex
var subscribers = make(map[event][]Callback) // key: either an ID or "SubmissionCreated"
var store = make(map[xid.ID]Submission)

func SubscribeSubmission(id string, cb Callback) {
	sid, e := xid.FromString(id)
	if e != nil {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	if s := getSubmission(sid); s != nil {
		cb(s)
		subscribe(submissionChanged{sid}, cb)
	}
}

func SubscribeCreations(cb Callback) {
	mutex.Lock()
	defer mutex.Unlock()
	subscribe(submissionCreated{}, cb)
}

func subscribe(e event, cb Callback) {
	subscribers[e] = append(subscribers[e], cb)
}

func CreateSubmission(q []byte) *Submission {
	mutex.Lock()
	defer mutex.Unlock()
	id := xid.New()
	s := &Submission{
		ID:         id,
		LastUpdate: timestamp(),
	}
	s.Question = append(s.Question, q...)

	store[id] = *s
	handleEvent(submissionCreated{}, s)
	return s
}

func UpdateSubmission(id string, a []byte) error {
	sid, e := xid.FromString(id)
	if e != nil {
		return fmt.Errorf("Invalid id: %s", id)
	}

	mutex.Lock()
	defer mutex.Unlock()

	s := getSubmission(sid)
	if s == nil {
		return fmt.Errorf("Submission not found")
	}

	if s.Answer != nil {
		return fmt.Errorf("Cannot modify a completed submission")
	}

	s.Answer = append(s.Answer, a...)
	s.LastUpdate = timestamp()
	handleEvent(submissionChanged{sid}, s)
	return nil
}

func getSubmission(id xid.ID) *Submission {
	if s, e := store[id]; e {
		return &s
	}

	return nil
}

func handleEvent(e event, s *Submission) {
	log.Info("Posting event: %s", e)

	if cbs, exists := subscribers[e]; exists {
		for _, cb := range cbs {
			go cb(s)
		}
	}
}

func timestamp() int64 {
	return time.Now().Unix()
}
