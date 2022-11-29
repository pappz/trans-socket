package test_helper

import (
	"bufio"
	"net"
)

type ServerClient struct {
	Conn          net.Conn
	ConnServer    net.Conn
	listener      net.Listener
	connReadyFlag chan bool
	alive         bool
}

func NewServerClient() (*ServerClient, error) {
	var err error
	mc := &ServerClient{
		connReadyFlag: make(chan bool),
		alive:         true,
	}

	mc.listener, err = net.Listen("tcp", ":1234")
	if err != nil {
		return nil, err
	}

	var acceptErr error
	go func() {
		mc.ConnServer, acceptErr = mc.listener.Accept()
		mc.connReadyFlag <- true
	}()

	mc.Conn, err = net.Dial("tcp", ":1234")
	if err != nil {
		return nil, err
	}

	<-mc.connReadyFlag
	if acceptErr != nil {
		return nil, acceptErr
	}

	return mc, nil
}

func (mc *ServerClient) ReadLine() (string, error) {
	connBuf := bufio.NewReader(mc.ConnServer)
	str, err := connBuf.ReadString('\n')
	return str, err
}

func (mc *ServerClient) SendSampleLine() error {
	_, err := mc.Conn.Write([]byte("Are you there?\n"))
	return err
}

func (mc *ServerClient) Teardown() {
	_ = mc.listener.Close()
}
