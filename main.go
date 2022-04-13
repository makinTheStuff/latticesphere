package main

import (
	// _ "net/http/pprof"

	"latticesphere/handlers"
	"latticesphere/networking/broadcaster"

	"fmt"
	"net/http"
	//"os"
	//"os/signal"
)

func main() {

	b := broadcaster.NewBoradcaster()

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/broadcast", b.WSProxy)
	http.HandleFunc("/bstruct", b.Bstruct)
	//http.HandleFunc("/path", rootHandler)

	go http.ListenAndServe(":8888", nil)
	fmt.Println("Visit http://0.0.0.0:8888")

	b.Run()

	// b.Wait()

	//sigCh := make(chan os.Signal)
	//signal.Notify(sigCh, os.Interrupt)
	//<-sigCh
	//signal.Stop(sigCh)
	//signal.Reset(os.Interrupt)
	//server.Shutdown()
}
