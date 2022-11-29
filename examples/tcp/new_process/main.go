package main

import (
	"bufio"
	"log"
	"net"

	"github.com/pappz/trans-socket/transsocket"
)

func readData(c net.Conn) (string, error) {
	connBuf := bufio.NewReader(c)
	return connBuf.ReadString('\n')

}

func onReceivedClient(conn net.Conn) {
	for {
		s, err := readData(conn)
		if err != nil {
			log.Printf("read err: %s", err.Error())
			break
		}

		log.Printf("msg from client: %s", s)
	}
}

func main() {
	log.Println("hello I am the new process!")

	s, err := transsocket.NewReceiver(0)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	log.Println("wait for transocket service connection")
	if err := s.WaitForSender(); err != nil {
		log.Fatal(err)
	}

	_, fd, err := s.RecvFileDescriptor()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("on new connection")
	onReceivedClient(fd)
}
