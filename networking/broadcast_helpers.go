package networking

import (
	"strconv"
)

func (b *Broadcaster) SubCountInt() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.subscribers)
}

func (b *Broadcaster) SubCountString() string {
	return strconv.Itoa(b.SubCountInt())
}

func (b *Broadcaster) SubCountB() []byte {
	return []byte(b.SubCountString())
}

func (b *Broadcaster) IsRunning() bool {
	b.RLock()
	defer b.RUnlock()

	return b.running
}

func (b *Broadcaster) SubscriberIDs() []ID {
	b.RLock()
	defer b.RUnlock()
	ids := make([]ID, len(b.subscribers))
	index := 0
	for id, _ := range b.subscribers {
		ids[index] = id
		index++
	}
	return ids
}

func (b *Broadcaster) Start() {
	// ::todo:: turn on and off for times when no new connections are made
	// and have add trigger this
	//
	// maybe also have a queue manager that drains and monitors queues
	// for different purposes (write, read, check, schedule, etc)
	// fmt.Println(1111, b.IsRunning(), b.SubCountInt())

	// if !b.IsRunning() && b.SubCountInt() > 0 {
	if !b.IsRunning() {
		// this check is to ensure the process has not already
		// been started or the sub we added that triggered
		// starting the service hasnt already been removed and
		// now the pool is empty (one user trying to reconnnect)
		// fmt.Println(1111)
		b.start()
	}
}
