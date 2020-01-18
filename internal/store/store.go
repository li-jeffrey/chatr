package store

import (
	"chatr/internal/logger"
	"fmt"
	"sync"

	"github.com/rs/xid"
)

var log = logger.GetLogger("store")

// Submission is a record containing a question-answer pair and its ID.
type Submission struct {
	ID         xid.ID
	Question   []byte
	Answer     []byte
	LastUpdate int64
	OwnerID    xid.ID
	AssignedID xid.ID
}

var Assigner func() xid.ID

type Event struct {
	Type       EventType
	Submission Submission
}

type EventType string

const (
	Assigned  EventType = "Assigned"
	Submitted EventType = "Submitted"
)

// stored state
var mutex sync.Mutex
var subscribers = make(map[xid.ID]chan<- Event)

func Subscribe(sessionID xid.ID) <-chan Event {
	events := make(chan Event)
	mutex.Lock()
	defer mutex.Unlock()
	subscribers[sessionID] = events

	return events
}

func CreateSubmission(q []byte, ownerID []byte) Submission {
	var oid xid.ID
	oid.UnmarshalText(ownerID)

	s := Submission{
		OwnerID: oid,
	}

	s.Question = append(s.Question, q...)

	inserted := insert(s)
	subscribers[oid] <- Event{
		Submitted, inserted,
	}

	go assign(inserted)
	return inserted
}

func UpdateSubmission(id string, a []byte) error {
	subID, _ := xid.FromString(id)
	s, e := getByID(subID)
	if !e {
		return fmt.Errorf("Submission not found")
	}

	if s.Answer != nil {
		return fmt.Errorf("Cannot modify a completed submission")
	}

	s.Answer = append(s.Answer, a...)

	insert(s)
	go sendUpdate(s)
	return nil
}

func assign(s Submission) {
	id := Assigner()
	if id == s.OwnerID {
		id = Assigner()
	}

	s.AssignedID = id
	inserted := insert(s)
	subscribers[id] <- Event{
		Assigned, inserted,
	}
}

func sendUpdate(s Submission) {
	if owner, e := subscribers[s.OwnerID]; e {
		owner <- Event{
			Submitted, s,
		}
	}

	if assigned, e := subscribers[s.AssignedID]; e {
		assigned <- Event{
			Assigned, s,
		}
	}
}
