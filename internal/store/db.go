package store

import (
	"sync"
	"time"

	"github.com/rs/xid"
)

var dbMutex sync.Mutex
var db = make(map[xid.ID] /* submissionId */ Submission)
var ownerIndex = make(map[xid.ID] /* ownerId */ map[xid.ID]struct{} /* submissionId */)
var assignedIndex = make(map[xid.ID] /* assignedId */ map[xid.ID]struct{} /*submissionId*/)

func insert(sub Submission) Submission {
	id := sub.ID
	isNew := id.IsNil()

	if isNew {
		id = xid.New()
		sub.ID = id
	}

	sub.LastUpdate = timestamp()
	dbMutex.Lock()
	defer dbMutex.Unlock()

	db[id] = sub

	ownerID := sub.OwnerID
	idx := ownerIndex[ownerID]
	if idx == nil {
		ownerIndex[ownerID] = make(map[xid.ID]struct{})
	}

	ownerIndex[ownerID][id] = struct{}{}

	if assignedID := sub.AssignedID; !assignedID.IsNil() {
		idx := assignedIndex[assignedID]
		if idx == nil {
			assignedIndex[assignedID] = make(map[xid.ID]struct{})
		}

		assignedIndex[assignedID][id] = struct{}{}
	}

	return sub
}

func getByID(id xid.ID) (Submission, bool) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	s, e := db[id]
	return s, e
}

func GetAllByOwnerID(ownerID xid.ID) []Submission {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	ids := ownerIndex[ownerID]
	if ids == nil {
		return nil
	}

	i := 0
	subs := make([]Submission, len(ids))
	for id := range ids {
		subs[i] = db[id]
		i++
	}

	return subs
}

func GetAllByOwnerIDOrderBy(ownerID xid.ID, ordering Ordering) []Submission {
	subs := GetAllByOwnerID(ownerID)
	ordering.by(subs)

	return subs
}

func GetAllByAssignedID(assignedID xid.ID) []Submission {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	ids := assignedIndex[assignedID]
	if ids == nil {
		return nil
	}

	i := 0
	subs := make([]Submission, len(ids))
	for id := range ids {
		subs[i] = db[id]
		i++
	}

	return subs
}

func GetAllByAssignedIDOrderBy(assignedID xid.ID, ordering Ordering) []Submission {
	subs := GetAllByAssignedID(assignedID)
	ordering.by(subs)

	return subs
}

func timestamp() int64 {
	return time.Now().Unix()
}
