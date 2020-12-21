package subscriber

import (
	"fmt"
	"github.com/gobwas/ws"
	"time"

	ln "latticesphere/networking"
	lm "latticesphere/networking/messages"
)

func (sub *Subscriber) Start() {
	sub.start()
}

func (sub *Subscriber) HasExpired() bool {
	sub.RLock()
	defer sub.RUnlock()
	return sub.expired
}

func (sub *Subscriber) ID() ln.ID {
	sub.RLock()
	defer sub.RUnlock()
	return sub.id
}

func (sub *Subscriber) ReadHeader() (ws.Header, error) {
	sub.Lock()
	defer sub.Unlock()
	sub.conn.SetDeadline(
		time.Now().Add(time.Duration(100) * time.Millisecond),
	)
	h, err := ws.ReadHeader(sub.conn)
	sub.conn.SetDeadline(
		time.Now().Add(time.Duration(2) * time.Millisecond),
	)
	return h, err
}

func (sub *Subscriber) GetRemoteAddr() string {
	sub.RLock()
	defer sub.RUnlock()
	return sub.conn.RemoteAddr().String()
}

func (sub *Subscriber) Send(msg lm.MsgCoordinate) bool {
	fmt.Println("Subscriber.Send", sub.ID(), sub.expired, msg)
	sub.Lock()
	defer sub.Unlock()
	// if !sub.expired && !msg.IsInTransport(sub.id) {
	sub.outgoing <- msg
	return true
	// }
	// return false
}
