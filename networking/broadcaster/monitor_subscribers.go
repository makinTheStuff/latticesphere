package broadcaster

import (
	"fmt"
	"time"
)

func (b *Broadcaster) Monitor() {
	go b.drainMessages()
	go b.monitorSubs()
}

func (b *Broadcaster) drainMessages() {
	for mc := range b.messages.MQueue {
		b.RLock()
		sub, found := b.subscribers[mc.RecipientID]
		b.RUnlock()
		if found {
			fmt.Println("drainmessages ------ ", mc, b.SubscriberIDs())
			fmt.Println("-----------------", sub.ID(), found)
			sub.Send(mc)
		}
	}
}

func (b *Broadcaster) monitorSubs() {
	fmt.Println("monitor")
	for {
		b.RLock()
		for _, sub := range b.subscribers {
			if sub.HasExpired() {
				// if the subscirber has already expired then
				// the connection has beenn closed and deleted
				// so its safe to remove
				b.RUnlock()
				fmt.Printf("\n\n---------------------expired %v %v %v %v", sub.ID(), b.SubscriberIDs(), b.lastID, len(b.subscribers))
				b.Remove(sub.ID())
				fmt.Printf("\n\n---------------------expired %v %v %v %v", sub.ID(), b.SubscriberIDs(), b.lastID, len(b.subscribers))
				b.RLock()
			}

		}
		b.RUnlock()
		// if b.SubCountInt() == 0 && !b.IsRunning() {
		//	b.stop()
		//	break
		// }
		// cool down (maybe diff strategy?)
		time.Sleep(2 * time.Millisecond)
	}
	fmt.Println("monitor")
}
