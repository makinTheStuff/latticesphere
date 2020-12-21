package subscriber

import (
	"net"
	"sync"

	ln "latticesphere/networking"
	lm "latticesphere/networking/messages"
)

type Subscriber struct {
	id        ln.ID                 `json:"id"`
	conn      net.Conn              `json:"-"`
	running   bool                  `json:"running"`
	expired   bool                  `json:"expired"`
	outgoing  chan lm.MsgCoordinate `json:"-"`
	messages  *lm.MessageContainer  `json:"-"` // to grab message; m.String()
	frameRate int                   `json:"frame_rate"`

	// createdAt, SDeletedOn, scheduled []

	sync.RWMutex `json:"-"`
}

func NewSubscriber(id ln.ID, conn net.Conn, frameRate int, messages *lm.MessageContainer) *Subscriber {
	return &Subscriber{
		id:        id,
		conn:      conn,
		frameRate: frameRate,
		running:   false, // the subscriber will start this
		expired:   false,
		outgoing:  make(chan lm.MsgCoordinate),
		messages:  messages,
	}
}
