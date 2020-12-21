package broadcaster

import (
	"fmt"
	"log"
	"net"

	ln "latticesphere/networking"
	lm "latticesphere/networking/messages"
	ls "latticesphere/networking/subscriber"
)

func (b *Broadcaster) start() {
	fmt.Println("\nBroadcaster.start  -- start")
	b.Lock()
	b.running = true
	b.Unlock()
	b.Monitor()
	fmt.Println("\nBroadcaster.start  -- start")
}

func (b *Broadcaster) stop() {
	fmt.Println("Broadcaster.stop -- start")
	b.Lock()
	defer b.Unlock()
	b.running = false
	fmt.Println("Broadcaster.stop -- end")
}

func (b *Broadcaster) AddSub(c net.Conn) {
	fmt.Println("\n\nBroadcaster.AddSub -- start")
	b.Lock()
	if b.subscribers == nil {
		fmt.Println("\tno subs found; making new map")
		b.subscribers = make(map[ln.ID]*ls.Subscriber)
	}
	b.Unlock()

	b.Start()
	fmt.Println(b.IsRunning())

	b.Lock()
	b.lastID += 1
	sub := ls.NewSubscriber(b.lastID, c, 2, &b.messages)
	b.subscribers[sub.ID()] = sub
	b.Unlock()

	sub.Start()
	fmt.Println("sub ids", b.SubscriberIDs())
	fmt.Println("Broadcaster.AddSub -- end\n")
}

func (b *Broadcaster) Remove(sid ln.ID) {
	fmt.Println("\tbefore: ", b.SubscriberIDs())
	b.Lock()
	delete(b.subscribers, sid)
	b.Unlock()
	fmt.Println("\t after: ", b.SubscriberIDs(), "\n\n")
}

func (b *Broadcaster) QueueMessage(msg *lm.Message) {
	b.messages.AddMessage(msg)
}

func (b *Broadcaster) upgradeToWS(conn net.Conn) {
	fmt.Println("\n\n")
	fmt.Println("start upgradeToWS\n")
	_, err := WS_UPGRADER.Upgrade(conn)
	if err != nil {
		fmt.Println("\tupgrade error: %s", err)
	} else {
		fmt.Println("\tadding sub")
		b.AddSub(conn)
		fmt.Println("\tnew sub ids", b.SubscriberIDs())

		// let the people know (just to test)
		// msg := fmt.Sprintf("{\"num_open_connns\": \"%[1]d\"}", b.SubCountInt())
		msg_template := "\nnum_open_connns: %[1]d; remote_addr: %[2]s"
		msg := fmt.Sprintf(msg_template, b.SubCountInt(), conn.RemoteAddr().String())
		fmt.Println(msg)
	}
	fmt.Println("end upgradeToWS\n")
}

func (b *Broadcaster) Run() {
	ln, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		b.upgradeToWS(conn)
	}
}
