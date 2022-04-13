package handlers

import (
	//	"fmt"
	"io"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("here")
	io.WriteString(w, `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8"/>
    <title>Sample of websocket with Golang</title>
  </head>
  <body>
	<div id="text"></div>
	<input type="text" id="text-view" /><br/>
    <input type="Button" id="send-button" value="Send" ></input>
    <script>
      var ws = new WebSocket("ws://" + document.domain + ":8889/ws");
	  ws.binaryType = 'arraybuffer';
      ws.onmessage = function(e) {
		var d = document.createElement("div");
        d.innerHTML = JSON.stringify(e.data);
        document.getElementById("text").appendChild(d);
		// ws.send("sdlfkhalsdkfjh");
		console.log(e)
		console.log(ws)
		
		// console.log(JSON.stringify(e.data.text))
		// create a JSON object
		// var jsonObject = JSON.parse(e.data);
		// var num_open_connns = jsonObject.num_open_connns;
		// console.log(num_open_connns)
      }
	  var textView = document.getElementById("text-view");
	  var buttonSend = document.getElementById("send-button");
		// Handle the button click event.
		buttonSend.onclick = function() {
		 // Send the data if connection is open
		  if (ws.readyState === WebSocket.OPEN){
			ws.send(textView.value);
			textView.value = "";
		  }
		}
	  ws.onclose = function(e){
		var d = document.createElement("div");
		d.innerHTML = "CLOSED";
        document.getElementById("text").appendChild(d);
	  }
    </script>
  </body>
</html>`)
}
