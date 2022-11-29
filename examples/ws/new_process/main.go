package main

import (
	"log"
	"net"

	"github.com/webkeydev/websocket"

	"github.com/pappz/trans-scoket/examples/wsgenerator"
	"github.com/pappz/trans-scoket/transsocket"
)

func fdToWsConn(conn net.Conn) *websocket.Conn {
	wsc := wsgenerator.NewWsConn()
	wsc.SetUnderlyingConn(conn)
	return wsc
}

func main() {
	log.Println("Hello, I am the new process!")

	log.Println("wait for socket sender")
	s, _ := transsocket.NewReceiver(0)
	if err := s.WaitForSender(); err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// wait for file descriptor
	_, fd, err := s.RecvFileDescriptor()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("on new file descriptor")

	// generate websocket connection
	wsConn := fdToWsConn(fd)

	// try out the connection
	log.Println("ready to use the ws connection")
	mt, message, err := wsConn.ReadMessage()
	if err != nil {
		log.Fatalf("failed to read msg: %s", err)
	}
	log.Printf("received msg: %s", message)

	err = wsConn.WriteMessage(mt, message)
	if err != nil {
		log.Fatalf("failed to read msg: %s", err)
	}
	log.Printf("sent msg: %s", message)
}
