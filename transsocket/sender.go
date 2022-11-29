package transsocket

import (
	"errors"
	"net"
	"syscall"
)

var (
	errMismatchLen    = errors.New("length of write data is mismatched")
	errMismatchOOBLen = errors.New("length of oob data is mismatched")
	errTooLargeData   = errors.New("too large data")
)

type Sender struct {
	conn *net.UnixConn
}

func NewSender() Sender {
	return Sender{}
}

func (s *Sender) Connect() error {
	c, err := net.Dial("unix", sockAddr)
	if err != nil {
		return err
	}
	s.conn = c.(*net.UnixConn)
	return nil
}

func (s *Sender) Disconnect() error {
	return s.conn.Close()
}

func (s *Sender) SendTCPFileDescriptor(c net.Conn) error {
	return s.send(nil, c)
}

func (s *Sender) SendTCPFileDescriptorWithData(data []byte, c net.Conn) error {
	if len(data) > dataLimit {
		return errTooLargeData
	}

	return s.send(data, c)
}

func (s *Sender) send(d []byte, c net.Conn) error {
	f, err := c.(*net.TCPConn).File()
	if err != nil {
		return err
	}

	rights := syscall.UnixRights(int(f.Fd()))
	n, oobn, err := s.conn.WriteMsgUnix(d, rights, nil)
	if err != nil {
		return err
	}

	if n != len(d) {
		return errMismatchLen
	}

	if oobn != len(rights) {
		return errMismatchOOBLen
	}
	_ = f.Close()
	return nil
}
