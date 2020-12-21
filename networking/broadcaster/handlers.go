package broadcaster

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	lm "latticesphere/networking/messages"
)

func (b *Broadcaster) WSProxy(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("here")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	msg := lm.NewMessage(string(reqBody), 1, b.SubscriberIDs())
	// for _, id := range b.SubscriberIDs() {
	//	fmt.Println("rid: ", id)
	//	msg.AddRecipient(id)
	// }
	fmt.Println(msg)
	b.QueueMessage(msg)
	io.WriteString(w, "Sent")
}

func (b *Broadcaster) Bstruct(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("here")
	b.Lock()
	fmt.Println(b)
	mapB, err := json.Marshal(b)
	fmt.Println(err)
	io.WriteString(w, string(mapB))
	fmt.Println(mapB)
	b.Unlock()
}
