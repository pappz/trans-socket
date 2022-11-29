package main

import (
	"bufio"
	"log"
	"net"

	"github.com/pappz/trans-socket/transsocket"
)

func readData(c *net.Conn) (string, error) {
	connBuf := bufio.NewReader(*c)
	return connBuf.ReadString('\n')
}

func main() {
	log.Println("Hello, I am the old process")
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Println("wait to test tcp client")
	var conn net.Conn
	if conn, err = l.Accept(); err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("on new tcp client")

	s, err := readData(&conn)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("received data: %s", s)

	c := transsocket.NewSender()

	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := c.SendTCPFileDescriptor(conn); err != nil {
		log.Fatalf("%v", err)
	}

	_ = conn.Close()
}
