package sessions

import (
	"bytes"
	"chatr/internal/store"

	"github.com/valyala/fastjson"
)

var pool fastjson.ArenaPool

// impl for sessions
var (
	submission []byte = []byte("Submission")
	assignment []byte = []byte("Assignment")
	ping       []byte = []byte("Ping")
	pong       []byte = []byte("Pong")
)

var byLatest store.Ordering = func(s1, s2 *store.Submission) bool {
	return s1.LastUpdate > s2.LastUpdate
}

var byEarliest store.Ordering = byLatest.Reversed()

func (s *Session) open() {
	s.sendConnect()
	for _, sub := range store.GetAllByOwnerIDOrderBy(s.ID, byEarliest) {
		s.sendResponse(submission, sub)
	}

	for _, assign := range store.GetAllByAssignedIDOrderBy(s.ID, byLatest) {
		s.sendResponse(assignment, assign)
	}

	events := store.Subscribe(s.ID)

Listener:
	for {
		select {
		case e := <-events:
			switch e.Type {
			case store.Submitted:
				s.sendResponse(submission, e.Submission)
			case store.Assigned:
				s.sendResponse(assignment, e.Submission)
			}
		case msg, open := <-s.Reads:
			if !open {
				break Listener
			}

			if v, e := s.p.ParseBytes(msg); e == nil {
				if bytes.Equal(ping, v.GetStringBytes("RequestType")) {
					a := pool.Get()
					a.Reset()
					obj := a.NewObject()
					obj.Set("ResponseType", a.NewStringBytes(pong))
					s.sendMessage(obj.MarshalTo(nil))
				}
			}
		case <-s.closeSignal:
			break Listener
		}
	}

	mutex.Lock()
	s.isActive = false
	mutex.Unlock()

	close(s.Writes)
	log.Info("Closed session %s.", s.ID)
}

func (s *Session) close() {
	s.closeSignal <- struct{}{}
	close(s.closeSignal)

	// wait until all channels are closed
	<-s.Writes
	<-s.Reads
}

func (s *Session) sendConnect() {
	a := pool.Get()
	a.Reset()
	obj := a.NewObject()
	obj.Set("ResponseType", a.NewString("Connection"))

	id, _ := s.ID.MarshalText()
	obj.Set("SessionID", a.NewStringBytes(id))

	s.sendMessage(obj.MarshalTo(nil))
}

func (s *Session) sendDisconnect(reason string) {
	a := pool.Get()
	a.Reset()
	obj := a.NewObject()
	obj.Set("ResponseType", a.NewString("Disconnection"))
	obj.Set("Reason", a.NewString(reason))

	s.sendMessage(obj.MarshalTo(nil))
}

func (s *Session) sendResponse(responseType []byte, sub store.Submission) {
	a := pool.Get()
	a.Reset()
	obj := a.NewObject()
	obj.Set("ResponseType", a.NewStringBytes(responseType))
	id, _ := sub.ID.MarshalText()
	obj.Set("ID", a.NewStringBytes(id))
	obj.Set("Question", a.NewStringBytes(sub.Question))
	obj.Set("Answer", a.NewStringBytes(sub.Answer))
	obj.Set("LastUpdate", a.NewNumberInt(int(sub.LastUpdate)))

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
