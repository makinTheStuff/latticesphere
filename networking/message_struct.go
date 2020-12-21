package networking

import (
	"sync"
)

type Message struct {
	// innclude header ?
	id         ID
	err        error
	content    string
	senderID   ID
	recipients map[ID]Status

	// sendAt int, createdAt, sentAt
	sync.RWMutex
}

func NewMessage(content string, sid ID, rids []ID) *Message {
	m := Message{
		content:    content,
		senderID:   sid,
		recipients: make(map[ID]Status),
	}
	m.AddRecipients(rids)
	return &m
}

type MsgCoordinate struct {
	messsageID  ID
	recipientID ID
}

func NewMsgCoordinate(mID, rID ID) MsgCoordinate {
	return MsgCoordinate{
		messsageID:  mID,
		recipientID: rID,
	}
}

type Status struct {
	// type and status consts defined in consts.go
	status MessageStatus
	err    error

	sync.RWMutex
}

func (s *Status) updateStatus(status MessageStatus, err error) {
	s.Lock()
	defer s.Unlock()
	s.status = status
	s.err = err
}
