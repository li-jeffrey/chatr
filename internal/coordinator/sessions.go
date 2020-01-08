package coordinator

import (
	"chatr/internal/logger"
	"chatr/internal/store"
	"sync"

	"github.com/rs/xid"
	"github.com/valyala/fastjson"
)

var log = logger.GetLogger("coordinator")

var pool fastjson.ArenaPool

var mutex sync.Mutex
var sessions = make(map[xid.ID]Session, 0)
var sessionIds []xid.ID // this is used as a fifo circular structure
var pendingAssignment []chan<- *Session

type Session struct {
	ID     xid.ID
	Reads  chan []byte
	Writes chan []byte
	p      fastjson.Parser
}

func CreateSession() *Session {
	mutex.Lock()
	defer mutex.Unlock()

	id := xid.New()
	s := &Session{
		ID:     id,
		Reads:  make(chan []byte),
		Writes: make(chan []byte),
	}
	sessions[id] = *s

	if len(pendingAssignment) != 0 {
		var c chan<- *Session
		c, pendingAssignment = pendingAssignment[0], pendingAssignment[1:]
		c <- s
	}

	sessionIds = append(sessionIds, id)
	go onSessionCreated(s)
	log.Info("Created session %s.", id)
	return s
}

func onSessionCreated(s *Session) {
	for msg := range s.Reads {
		if v, err := s.p.ParseBytes(msg); err == nil {
			handleMessage(s, v)
		} else {
			s.sendError(BadMessage, "Could not parse message")
		}
	}
	close(s.Writes)
	removeSession(s.ID)
	log.Info("Removed session %s.", s.ID)
}

func removeSession(id xid.ID) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(sessions, id)
}

func selectNextSession() <-chan *Session {
	c := make(chan *Session, 1)
	mutex.Lock()
	defer mutex.Unlock()
	return selectNextSessionRecur(c)
}

func selectNextSessionRecur(c chan *Session) <-chan *Session {
	// pop front of queue
	if len(sessionIds) == 0 {
		pendingAssignment = append(pendingAssignment, c)
		return c
	}

	var next xid.ID
	next, sessionIds = sessionIds[0], sessionIds[1:]
	s, e := sessions[next]
	if !e {
		// session removed: select next one
		return selectNextSessionRecur(c)
	}

	// put selected to the end of the queue
	sessionIds = append(sessionIds, next)
	go func() {
		c <- &s
		close(c)
	}()

	return c
}

func (s *Session) sendResponse(responseType []byte, sub *store.Submission) {
	a := pool.Get()
	a.Reset()
	obj := a.NewObject()
	obj.Set("ResponseType", a.NewStringBytes(responseType))
	id, _ := s.ID.MarshalText()
	obj.Set("ID", a.NewStringBytes(id))
	obj.Set("Question", a.NewStringBytes(sub.Question))
	obj.Set("Answer", a.NewStringBytes(sub.Answer))

	s.sendMessage(obj.MarshalTo(nil))
}

func (s *Session) sendError(errorType []byte, errorMsg string) {
	a := pool.Get()
	a.Reset()
	obj := a.NewObject()
	obj.Set("Error", a.NewStringBytes(errorType))
	obj.Set("Message", a.NewString(errorMsg))

	s.sendMessage(obj.MarshalTo(nil))
}

func (s *Session) sendMessage(b []byte) {
	s.Writes <- b
}
