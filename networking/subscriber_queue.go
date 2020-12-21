package networking

import (
	"fmt"
	// "io"
	// "net"
	//"time"

	"github.com/gobwas/ws"
)

func (sub *Subscriber) start() {
	sub.Lock()

	fmt.Println("\n\nSubscriber.start  -- start")
	sub.expired = false
	sub.running = true

	fmt.Println("fire up queues")
	go sub.drainOutgoingMessages()
	go sub.drainIncomingMessages()
	fmt.Println("\nSubscriber.start  --  end")

	sub.Unlock()
}

func (sub *Subscriber) drainIncomingMessages() {
	fmt.Println("\n\ndrainIncomingMessages -- start\n")
	for {
		// fmt.Println("\treading header")
		// blocking call to read header
		header, err := sub.ReadHeader()
		if err != nil {
			// handle error
			// fmt.Println("\tb.ReadHeader error: %s", err)
			sub.hadleConnErr(err)
		} else {
			if header.OpCode == ws.OpClose {
				fmt.Println("\n\tbefore: %+q", sub.ID())
				sub.deleteConn()
				fmt.Println("\tafter: %+q", sub.ID())
				break // continue
			}
			if !sub.HasExpired() {
				fmt.Println(
					"\tsub.callback -- sub.HasExpired: ",
					sub.HasExpired(), ", sub.ID: ", sub.ID(),
				)
				sub.callback(header)
			} else {
				// should be cleaned up by the broadcaster
				fmt.Println("\tsub.HasExpired", sub.HasExpired(), sub.ID())
				break
			}
		}
	}
	fmt.Println("\ndrainIncomingMessages -- end")
}

func (sub *Subscriber) drainOutgoingMessages() {
	fmt.Println("\n\ndrainOutgoingMessages -- start")
	for mc := range sub.outgoing {
		msg := sub.messages.GetPendingMessage(mc)
		if !msg.IsInTransport(sub.ID()) {
			msg.SetInTransport(sub.ID())
			if err := sub.writePayload(msg.Bytes()); err != nil {
				fmt.Println("\nsub write error: %s", err, sub.ID(), msg.String())
				sub.messages.Failed(mc, err)
			} else {
				sub.messages.Delivered(mc)
			}
		}

		// cool down (maybe diff strategy?)
		// time.Sleep(time.Duration(sub.frameRate) * time.Millisecond)
	}
	fmt.Println("\ndrainOutgoingMessages -- end")
}
