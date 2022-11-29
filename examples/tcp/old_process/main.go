package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pappz/trans-scoket/transsocket"
)

func readData(c *net.Conn) (string, error) {
	connBuf := bufio.NewReader(*c)
	return connBuf.ReadString('\n')
}

func main() {
	fmt.Println("Hello, I am the old process")
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("wait to test tcp client")
	var conn net.Conn
	if conn, err = l.Accept(); err != nil {
		log.Fatalf("%v", err)
	}

	s, err := readData(&conn)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("received data: %s", s)

	c := transsocket.NewSender()

	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := c.SendTCPFileDescriptor(conn); err != nil {
		log.Fatalf("%v", err)
	}

	_ = conn.Close()

	// wait to check the results by human
	time.Sleep(1 * time.Minute)
}
