package store

import (
	"chatr/internal/logger"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Submission is a record containing a question-answer pair and its ID.
type Submission struct {
	ID         string
	Question   []byte
	Answer     []byte
	LastUpdate int64
}

// EventType is a string representing the type of the event
type EventType string

type Callback func(*Submission)

var log = logger.GetLogger("store")

// events
type event struct {
	eventType EventType
	id        string
}

const submissionChanged EventType = "SubmissionChanged"

var submissionCreated = event{"SubmissionCreated", ""}

// stored state
var mutex sync.Mutex
var subscribers = make(map[event][]Callback)
var store = make(map[string]Submission)

func SubscribeSubmission(id string, cb Callback) {
	mutex.Lock()
	defer mutex.Unlock()
	if s := getSubmission(id); s != nil {
		cb(s)
		subscribe(event{submissionChanged, id}, cb)
	}
}

func SubscribeCreations(cb Callback) {
	mutex.Lock()
	defer mutex.Unlock()
	subscribe(submissionCreated, cb)
}

func subscribe(e event, cb Callback) {
	subscribers[e] = append(subscribers[e], cb)
}

func CreateSubmission(q []byte) *Submission {
	mutex.Lock()
	defer mutex.Unlock()
	id := uuid.New().String()
	s := &Submission{
		ID:         id,
		LastUpdate: timestamp(),
	}
	s.Question = append(s.Question, q...)

	store[id] = *s
	handleEvent(submissionCreated, s)
	return s
}

func UpdateSubmission(id string, a []byte) error {
	mutex.Lock()
	defer mutex.Unlock()
	s := getSubmission(id)
	if s == nil {
		return fmt.Errorf("Submission not found")
	}

	if s.Answer != nil {
		return fmt.Errorf("Cannot modify a completed submission")
	}

	s.Answer = append(s.Answer, a...)
	s.LastUpdate = timestamp()
	handleEvent(event{submissionChanged, s.ID}, s)
	return nil
}

func getSubmission(id string) *Submission {
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
