package networking

import "sync"

type Broadcaster struct {
	subscribers map[ID]*Subscriber `json:"subscribers"`
	messages    MessageContainer   `json:"-"`
	running     bool               `json:"running"`
	msgCount    ID                 `json:"msgCount"`
	lastID      ID                 `json:"-"`

	sync.RWMutex `json:"-"`
}

func NewBoradcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: map[ID]*Subscriber{},
		messages:    NewMessageContainer(),
		running:     false,
		msgCount:    0,
		lastID:      0,
	}
}
