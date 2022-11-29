package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/pappz/trans-scoket/transsocket"
)

var (
	onReceived = make(chan bool)
)

func readData(c net.Conn) (string, error) {
	connBuf := bufio.NewReader(c)
	return connBuf.ReadString('\n')

}

func onReceivedClient(conn net.Conn) {
	fmt.Println("on new connection")
	for {
		s, err := readData(conn)
		if err != nil {
			log.Printf("%v", err)
			break
		}

		log.Printf("msg from client: %s", s)
	}
	onReceived <- true
}

func main() {
	fmt.Println("hello I am the new process!")

	s, err := transsocket.NewReceiver(0)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.WaitForSender(); err != nil {
		log.Fatal(err)
	}

	_, fd, err := s.RecvFileDescriptor()
	if err != nil {
		log.Fatal(err)
	}
	onReceivedClient(fd)

	<-onReceived
}
