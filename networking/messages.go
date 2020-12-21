package networking

import (
	// "common"
	"sync"
)

// var messages MessageContainer

type MessageContainer struct {
	// maps of mid to msg
	delivered map[ID]*Message
	pending   map[ID]*Message
	failed    map[ID]*Message
	mqueue    chan MsgCoordinate // chan of incomig mids

	lastMsgID ID

	// sync.WaitGroup
	sync.RWMutex
}

func NewMessageContainer() MessageContainer {
	return MessageContainer{
		delivered: map[ID]*Message{},
		pending:   map[ID]*Message{},
		failed:    map[ID]*Message{},
		mqueue:    make(chan MsgCoordinate),
		lastMsgID: 0,
	}
}

func (mc *MessageContainer) GetPendingMessage(coord MsgCoordinate) *Message {
	mc.RLock()
	defer mc.RUnlock()
	return mc.pending[coord.messsageID]
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
		mc.mqueue <- NewMsgCoordinate(msg.id, rid)
	}
	msg.RUnlock()
}

func (mc *MessageContainer) Delivered(coord MsgCoordinate) {
	status := mc.pending[coord.messsageID].Delivered(coord.recipientID)

	mc.Lock()
	defer mc.Unlock()
	if status == COMPLETED {
		mc.delivered[coord.messsageID] = mc.pending[coord.messsageID]
		delete(mc.pending, coord.messsageID)
	}
}

func (mc *MessageContainer) Failed(coord MsgCoordinate, err error) {
	status := mc.pending[coord.messsageID].Failed(coord.recipientID, err)

	mc.Lock()
	defer mc.Unlock()
	if status == COMPLETED {
		mc.failed[coord.messsageID] = mc.pending[coord.messsageID]
		delete(mc.pending, coord.messsageID)
	}
}
