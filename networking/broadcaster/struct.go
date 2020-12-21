package broadcaster

import (
	"sync"

	ln "latticesphere/networking"
	lm "latticesphere/networking/messages"
	ls "latticesphere/networking/subscriber"
)

type Broadcaster struct {
	subscribers map[ln.ID]*ls.Subscriber `json:"subscribers"`
	messages    lm.MessageContainer      `json:"-"`
	running     bool                     `json:"running"`
	msgCount    ln.ID                    `json:"msgCount"`
	lastID      ln.ID                    `json:"-"`

	sync.RWMutex `json:"-"`
}

func NewBoradcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: map[ln.ID]*ls.Subscriber{},
		messages:    lm.NewMessageContainer(),
		running:     false,
		msgCount:    0,
		lastID:      0,
	}
}
