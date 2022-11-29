package main

import (
	"log"
	"net"

	"github.com/webkeydev/websocket"

	"github.com/pappz/trans-socket/examples/ws/generator"
	"github.com/pappz/trans-socket/transsocket"
)

func fdToWsConn(conn net.Conn) (*websocket.Conn, error) {
	wsg := generator.NewWsGenerator()
	wsc, err := wsg.NewWsConn()
	if err != nil {
		return nil, err
	}
	wsc.SetUnderlyingConn(conn)
	return wsc, nil
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
	wsConn, err := fdToWsConn(fd)
	if err != nil {
		log.Fatalf("failed to generate ws connection from fd: %s", err)
	}

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
