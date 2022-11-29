package main

import (
	"log"
	"net/http"

	"github.com/webkeydev/websocket"

	"github.com/pappz/trans-scoket/transsocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func testConn(wsc *websocket.Conn) error {
	mt, message, err := wsc.ReadMessage()
	if err != nil {
		return nil
	}
	log.Printf("recv: %s", message)

	err = wsc.WriteMessage(mt, message)
	if err != nil {
		return err
	}
	return nil
}

func transfer(w http.ResponseWriter, r *http.Request) {
	wsc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func(wsc *websocket.Conn) {
		_ = wsc.Close()
	}(wsc)

	if err := testConn(wsc); err != nil {
		log.Fatal(err)
	}

	tc := transsocket.NewSender()
	if err := tc.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := tc.SendTCPFileDescriptor(wsc.UnderlyingConn()); err != nil {
		log.Fatalf("%v", err)
	}
}

func main() {
	log.Print("Hello, I am the old process!")
	http.HandleFunc("/", transfer)
	log.Fatal(http.ListenAndServe(":1234", nil))
}
