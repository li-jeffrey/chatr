package store

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"sync"

	"github.com/rs/xid"
)

var fmutex sync.Mutex
var handle *os.File
var buf = bytes.NewBuffer(nil)
var us uint8 = 0x1f
var rs uint8 = 0x1e

func EnablePersistence(location string) {
	if f, e := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644); e != nil {
		log.Fatal("Failed to open %s: %s", location, e)
	} else {
		loadFromFile(f)
		handle = f
	}
}

func writeToLog(s Submission) {
	if e := doWrite(s); e != nil {
		log.Error("Failed to write: %s", e)
	}
}

func doWrite(s Submission) error {
	fmutex.Lock()
	defer fmutex.Unlock()
	buf.Reset()
	buf.Write(s.ID.Bytes())
	buf.WriteByte(us)
	buf.Write(s.Question)
	buf.WriteByte(us)
	buf.Write(s.Answer)
	buf.WriteByte(us)
	buf.Write(s.OwnerID.Bytes())
	buf.WriteByte(us)
	buf.Write(s.AssignedID.Bytes())
	buf.WriteByte(us)

	b := make([]byte, 8)
	binary.PutVarint(b, s.LastUpdate)
	buf.Write(b)
	buf.WriteByte(rs)

	_, e := buf.WriteTo(handle)
	return e
}

func loadFromFile(f *os.File) {
	loaded := make(map[xid.ID]Submission)
	scanner := bufio.NewScanner(f)
	scanner.Split(scanRecord)
	for scanner.Scan() {
		var s Submission
		record := scanner.Bytes()
		split := bytes.Split(record, []byte{us})
		if len(split) < 6 {
			log.Warn("Deformed record: %s", record)
			continue
		}

		id, _ := xid.FromBytes(split[0])
		ownerID, _ := xid.FromBytes(split[3])
		assignedID, _ := xid.FromBytes(split[4])
		lastUpdate, _ := binary.Varint(split[5])

		s.ID = id
		s.Question = split[1]
		s.Answer = split[2]
		s.OwnerID = ownerID
		s.AssignedID = assignedID
		s.LastUpdate = lastUpdate

		loaded[id] = s
	}

	db = loaded
	log.Info("Loaded %v entries from file.", len(loaded))
}

func scanRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, rs); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
