package transsocket

import (
	"errors"
	"net"
	"os"
	"syscall"
)

var (
	errInvalidNumberOfDescriptors = errors.New("invalid number of descriptors")
)

type Receiver struct {
	conn     *net.UnixConn
	listener net.Listener
	buffSize int
}

// NewReceiver The new instance will listen on UDN
func NewReceiver(buffSize int) (*Receiver, error) {
	if buffSize <= 0 {
		buffSize = 32
	}

	s := &Receiver{
		buffSize: buffSize,
	}

	return s, s.listen()
}

func (r *Receiver) WaitForSender() error {
	conn, err := r.listener.Accept()
	if err != nil {
		return err
	}

	r.conn = conn.(*net.UnixConn)
	return nil
}

func (r *Receiver) RecvFileDescriptor() ([]byte, net.Conn, error) {
	var receivedData []byte
	buf := make([]byte, dataLimit)
	oob := make([]byte, 32)
	bufLen, oobLen, _, _, err := r.conn.ReadMsgUnix(buf, oob)
	if err != nil {
		return receivedData, nil, err
	}

	receivedData = getDataFromBuffer(buf, bufLen)

	scms, err := syscall.ParseSocketControlMessage(oob[:oobLen])
	if err != nil {
		return receivedData, nil, err
	}

	if len(scms) != 1 {
		return receivedData, nil, errInvalidNumberOfDescriptors
	}

	scm := scms[0]
	fds, err := syscall.ParseUnixRights(&scm)
	if err != nil {
		return receivedData, nil, err
	}
	if len(fds) != 1 {
		return receivedData, nil, errInvalidNumberOfDescriptors
	}

	fd := os.NewFile(uintptr(fds[0]), "passed-fd")
	conn, err := net.FileConn(fd)
	_ = fd.Close()

	return receivedData, conn, err
}

func (r *Receiver) listen() (err error) {
	if err := os.RemoveAll(sockAddr); err != nil {
		return err
	}

	r.listener, err = net.Listen("unix", sockAddr)
	return
}

func getDataFromBuffer(buf []byte, bufLen int) []byte {
	if bufLen > 1 {
		return buf[:bufLen]
	}

	if buf[0] == 0 {
		return nil
	} else {
		return []byte{buf[0]}
	}
}
