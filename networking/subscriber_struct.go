package networking

import (
	"net"
	"sync"
)

type Subscriber struct {
	id        ID                 `json:"id"`
	conn      net.Conn           `json:"-"`
	running   bool               `json:"running"`
	expired   bool               `json:"expired"`
	outgoing  chan MsgCoordinate `json:"-"`
	messages  *MessageContainer  `json:"-"` // to grab message; m.String()
	frameRate int                `json:"frame_rate"`

	// createdAt, SDeletedOn, scheduled []

	sync.RWMutex `json:"-"`
}

func NewSubscriber(id ID, conn net.Conn, frameRate int, messages *MessageContainer) *Subscriber {
	return &Subscriber{
		id:        id,
		conn:      conn,
		frameRate: frameRate,
		running:   false, // the subscriber will start this
		expired:   false,
		outgoing:  make(chan MsgCoordinate),
		messages:  messages,
	}
}
