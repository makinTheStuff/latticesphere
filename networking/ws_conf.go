package networking

import (
	"net/http"
	//"os"
	//"os/signal"
	"runtime"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
)

// Prepare handshake header writer from http.Header mapping.
var HEADER = ws.HandshakeHeaderHTTP(http.Header{
	"X-Go-Version": []string{runtime.Version()},
})

var WS_UPGRADER = ws.Upgrader{
	OnHost: func(host []byte) error {
		//if string(host) == "github.com" {
		//	return nil
		//}
		//return ws.RejectConnectionError(
		//	ws.RejectionStatus(403),
		//	ws.RejectionHeader(ws.HandshakeHeaderString(
		//		"X-Want-Host: github.com\r\n",
		//	)),
		//)
		return nil
	},
	OnHeader: func(key, value []byte) error {
		if string(key) != "Cookie" {
			return nil
		}
		ok := httphead.ScanCookie(value, func(key, value []byte) bool {
			// Check session here or do some other stuff with cookies.
			// Maybe copy some values for future use.
			return true
		})
		if ok {
			return nil
		}
		return ws.RejectConnectionError(
			ws.RejectionReason("bad cookie"),
			ws.RejectionStatus(400),
		)
	},
	OnBeforeUpgrade: func() (ws.HandshakeHeader, error) {
		return HEADER, nil
	},
}
