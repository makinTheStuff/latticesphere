package messages

import (
	// "common"
	ln "latticesphere/networking"
	"sync"
)

type MessageContainer struct {
	// maps of mid to msg
	delivered map[ln.ID]*Message
	pending   map[ln.ID]*Message
	failed    map[ln.ID]*Message
	MQueue    chan MsgCoordinate // chan of incomig mids

	lastMsgID ln.ID

	sync.RWMutex
}

func NewMessageContainer() MessageContainer {
	return MessageContainer{
		delivered: map[ln.ID]*Message{},
		pending:   map[ln.ID]*Message{},
		failed:    map[ln.ID]*Message{},
		MQueue:    make(chan MsgCoordinate),
		lastMsgID: 0,
	}
}

func (mc *MessageContainer) GetPendingMessage(coord MsgCoordinate) *Message {
	mc.RLock()
	defer mc.RUnlock()
	return mc.pending[coord.MesssageID]
}

func (mc *MessageContainer) AddMessage(msg *Message) {
	mc.Lock()
	defer mc.Unlock()
	mc.lastMsgID += 1

	msg.Lock()
	msg.id = mc.lastMsgID
	msg.Unlock()

	mc.pending[mc.lastMsgID] = msg
	msg.RLock()
	for rid, _ := range msg.recipients {
		mc.MQueue <- NewMsgCoordinate(msg.id, rid)
	}
	msg.RUnlock()
}

func (mc *MessageContainer) Delivered(coord MsgCoordinate) {
	status := mc.pending[coord.MesssageID].Delivered(coord.RecipientID)

	mc.Lock()
	defer mc.Unlock()
	if status == ln.COMPLETED {
		mc.delivered[coord.MesssageID] = mc.pending[coord.MesssageID]
		delete(mc.pending, coord.MesssageID)
	}
}

func (mc *MessageContainer) Failed(coord MsgCoordinate, err error) {
	status := mc.pending[coord.MesssageID].Failed(coord.RecipientID, err)

	mc.Lock()
	defer mc.Unlock()
	if status == ln.COMPLETED {
		mc.failed[coord.MesssageID] = mc.pending[coord.MesssageID]
		delete(mc.pending, coord.MesssageID)
	}
}
