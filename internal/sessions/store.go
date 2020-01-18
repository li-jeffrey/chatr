package sessions

import (
	"chatr/internal/logger"
	"chatr/internal/store"
	"sync"

	"github.com/rs/xid"
	"github.com/valyala/fastjson"
)

var log = logger.GetLogger("sessions")

func init() {
	store.Assigner = func() xid.ID {
		s := <-selectNextSession()
		return s.ID
	}
}

var mutex sync.Mutex
var sessions = make(map[xid.ID]*Session, 0)
var sessionIds []xid.ID // used for assignment in a round-robin strategy
var pendingAssignment []chan<- *Session

type Session struct {
	ID          xid.ID
	Reads       chan []byte
	Writes      chan []byte
	isActive    bool
	p           fastjson.Parser
	closeSignal chan struct{}
}

func GetSession(sid []byte) *Session {
	mutex.Lock()
	defer mutex.Unlock()

	var id xid.ID
	if sid == nil {
		id = xid.New()
	} else {
		id.UnmarshalText(sid)
	}

	s, e := sessions[id]
	if !e {
		log.Info("Creating session %s.", id)
		s = &Session{ID: id}
		sessions[id] = s
		sessionIds = append(sessionIds, id)
	} else if s.isActive {
		log.Info("Closing duplicate session %s.", id)
		// disconnect the other session
		s.sendDisconnect("Duplicate session opened")
		s.close()
	} else {
		log.Info("Restoring session %s.", sid)
	}

	s.Reads = make(chan []byte)
	s.Writes = make(chan []byte)
	s.closeSignal = make(chan struct{}, 1)
	s.isActive = true

	assignIfPending(s)
	go s.open()

	return s
}

func assignIfPending(s *Session) {
	if len(pendingAssignment) > 0 {
		var c chan<- *Session
		c, pendingAssignment = pendingAssignment[0], pendingAssignment[1:]
		c <- s
	}
}

func selectNextSession() <-chan *Session {
	c := make(chan *Session, 1)
	mutex.Lock()
	defer mutex.Unlock()
	return selectNextSessionRecur(c)
}

func selectNextSessionRecur(c chan *Session) <-chan *Session {
	if len(sessionIds) < 2 {
		pendingAssignment = append(pendingAssignment, c)
		return c
	}

	var next xid.ID

	// shift queue
	next, sessionIds = sessionIds[0], sessionIds[1:]
	s := sessions[next]
	if !s.isActive {
		// session inactive: select next one
		return selectNextSessionRecur(c)
	}

	sessionIds = append(sessionIds, next)
	go func() {
		c <- s
		close(c)
	}()

	return c
}
