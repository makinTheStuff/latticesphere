package messages

import (
	ln "latticesphere/networking"
	"sync"
)

type Message struct {
	// innclude header ?
	id         ln.ID
	err        error
	content    string
	senderID   ln.ID
	recipients map[ln.ID]Status

	// sendAt int, createdAt, sentAt
	sync.RWMutex
}

func NewMessage(content string, sid ln.ID, rids []ln.ID) *Message {
	m := Message{
		content:    content,
		senderID:   sid,
		recipients: make(map[ln.ID]Status),
	}
	m.AddRecipients(rids)
	return &m
}

type MsgCoordinate struct {
	MesssageID  ln.ID
	RecipientID ln.ID
}

func NewMsgCoordinate(mID, rID ln.ID) MsgCoordinate {
	return MsgCoordinate{
		MesssageID:  mID,
		RecipientID: rID,
	}
}

type Status struct {
	// type and status consts defined in consts.go
	status ln.MessageStatus
	err    error

	sync.RWMutex
}

func (s *Status) updateStatus(status ln.MessageStatus, err error) {
	s.Lock()
	defer s.Unlock()
	s.status = status
	s.err = err
}
