package networking

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"syscall"
	"time"
	// "log"
	// "sort"
	// "strconv"

	"github.com/gobwas/ws"
)

func (sub *Subscriber) deleteConn() {
	fmt.Println("\n\nSubscriber.deleteConn -- start")
	sub.Lock()
	sub.expired = true
	sub.running = false
	sub.conn.Close()
	close(sub.outgoing)
	sub.Unlock()
	fmt.Println("\n\nSubscriber.deleteConn -- end")
}

func (sub *Subscriber) hadleConnErr(err error) error {
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, io.EOF) {
		sub.deleteConn()
		fmt.Println("sub.deleteConn -- error: %s", err)
	}
	return err
}

func (sub *Subscriber) writeFrame(frame ws.Frame) error {
	sub.Lock()
	defer sub.Unlock()
	fmt.Println("writeFrame")
	if err := ws.WriteFrame(sub.conn, frame); err != nil {
		fmt.Println("writeFrame: %s", err)
		return err
	}
	return nil
}

func (sub *Subscriber) writePayload(msgBody []byte) error {
	frame := ws.NewTextFrame(msgBody)
	fmt.Println("writePayload")
	if err := sub.hadleConnErr(sub.writeFrame(frame)); err != nil {
		fmt.Println("writePayload: %s", err)
		return err
	} else {
		fmt.Println("writeFrame")
		// sub.writeFrame(frame)
	}
	return nil
}

func (sub *Subscriber) readPayload(header ws.Header) []byte {
	// might want to use wsutil
	// https://godoc.org/github.com/gobwas/ws/wsutil#ControlHandler.HandleClose
	sub.Lock()
	fmt.Println("\n\nSubscriber.readPayload -- start")
	payload := make([]byte, header.Length)
	_, err := io.ReadFull(sub.conn, payload)
	sub.Unlock()

	if err != nil {
		fmt.Println("ws.ReadFrame error %s", err)
		sub.hadleConnErr(err)
		return []byte(err.Error())
	}
	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}
	fmt.Println("\n\nSubscriber.readPayload -- end")
	// sub.writePayload(payload)
	return payload
}

func (sub *Subscriber) callback(header ws.Header) {
	// set timeout and ther useful things in the future
	fmt.Println("\n\nSubscriber.callback -- start", header.Length)
	client := http.Client{Timeout: 1 * time.Second}
	payload := sub.readPayload(header)
	fmt.Println("payload: \t %s \n", payload)
	req, err := http.NewRequest("Post", "http://0.0.0.0:8888/broadcast", bytes.NewBuffer(payload))
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ioutil.ReadAll")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	//vfmt.Println("\nexpiring connection ", sub.ID())
	sb := string(body)
	fmt.Println("resp body:\t %s \n", sb)

	fmt.Println(reflect.TypeOf(payload).String())
	fmt.Println("payload:\t %s \n", string(payload))
	fmt.Println("\n\nSubscriber.callback -- end")
}
