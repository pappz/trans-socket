package main

import (
	"log"
	"net/http"

	"github.com/webkeydev/websocket"

	"github.com/pappz/trans-scoket/transsocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func testConn(wsc *websocket.Conn) error {
	mt, message, err := wsc.ReadMessage()
	if err != nil {
		return err
	}
	log.Printf("received msg: %s", message)

	err = wsc.WriteMessage(mt, message)
	if err != nil {
		return err
	}
	log.Printf("sent msg: %s", message)
	return nil
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("on new ws connection")
	wsc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("%s:", err)
		return
	}
	defer wsc.Close()

	err = testConn(wsc)
	if err != nil {
		log.Fatalf("%s", err)
	}

	tc := transsocket.NewSender()
	err = tc.Connect()
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer tc.Disconnect()

	if err := tc.SendTCPFileDescriptor(wsc.UnderlyingConn()); err != nil {
		log.Fatalf("%s", err)
	}
}

// Listen and wait to ws connection. If it happens transfer the connection to the new process
// For the test you can use the command below.
// wscat --connect ws://127.0.0.1:1234
func main() {
	log.Print("Hello, I am the old process!")
	http.HandleFunc("/", wsHandler)
	log.Fatal(http.ListenAndServe(":1234", nil))
}
