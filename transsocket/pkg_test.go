package transsocket

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pappz/trans-socket/transsocket/test_helper"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Test_transferFdWithData(t *testing.T) {
	wg := sync.WaitGroup{}

	sc, err := NewServerClient()
	assert.NoError(t, err)
	sampleMetaData := Meta{56}

	err = checkReadWrite(sc)
	assert.NoError(t, err)

	socketReceiver, err := NewReceiver(Size())
	defer socketReceiver.Close()
	assert.NoError(t, err)

	wg.Add(1)
	go func() {
		err := socketReceiver.WaitForSender()
		assert.NoError(t, err)

		metaData, netConn, err := socketReceiver.RecvFileDescriptor()
		assert.NoError(t, err)

		m, err := NewMetaFromBytes(metaData)
		assert.NoError(t, err)

		assert.Equal(t, sampleMetaData.DeviceId, m.DeviceId)

		w := bufio.NewWriter(netConn)
		_, err = w.Write([]byte("Are you there?\n"))
		assert.NoError(t, err)
		_ = w.Flush()

		_, err = sc.ReadLine()
		assert.NoError(t, err)
		wg.Done()
	}()

	c := NewSender()
	err = c.Connect()
	assert.NoError(t, err)

	data, err := sampleMetaData.ToSlice()
	assert.NoError(t, err)

	err = c.SendTCPFileDescriptorWithData(data, sc.Conn)
	assert.NoError(t, err)

	wg.Wait()
	sc.Teardown()
}

func Test_transferFd(t *testing.T) {
	wg := sync.WaitGroup{}
	sc, err := NewServerClient()
	assert.NoError(t, err)

	err = checkReadWrite(sc)
	assert.NoError(t, err)

	server, err := NewReceiver(Size())
	defer server.Close()
	assert.NoError(t, err)

	wg.Add(1)
	go func() {
		err := server.WaitForSender()
		assert.NoError(t, err)

		data, fd, err := server.RecvFileDescriptor()
		assert.NoError(t, err)

		assert.Nil(t, data)

		w := bufio.NewWriter(fd)
		_, err = w.Write([]byte("Are you there?\n"))
		assert.NoError(t, err)

		_ = w.Flush()

		_, err = sc.ReadLine()
		assert.NoError(t, err)

		wg.Done()

	}()

	c := NewSender()
	err = c.Connect()
	assert.NoError(t, err)

	err = c.SendTCPFileDescriptor(sc.Conn)
	assert.NoError(t, err)

	wg.Wait()
	sc.Teardown()
}

func Test_transferFullFilledBuffer(t *testing.T) {
	var netConn net.Conn
	wg := sync.WaitGroup{}
	sc, err := NewServerClient()
	assert.NoError(t, err)

	server, err := NewReceiver(Size())
	defer server.Close()
	assert.NoError(t, err)

	wg.Add(1)
	go func() {
		err := server.WaitForSender()
		assert.NoError(t, err)

		_, netConn, err = server.RecvFileDescriptor()
		assert.NoError(t, err)
		wg.Done()
	}()

	// fill buffer
	sizeOfRandString := 50
	payload := randStringBytes(sizeOfRandString)
	_, err = sc.ConnServer.Write(payload)
	assert.NoError(t, err)

	// read half
	readBuffer1 := make([]byte, sizeOfRandString/2)
	_, err = sc.Conn.Read(readBuffer1)
	assert.NoError(t, err)

	c := NewSender()
	err = c.Connect()
	assert.NoError(t, err)

	err = c.SendTCPFileDescriptor(sc.Conn)
	assert.NoError(t, err)

	wg.Wait()

	readBuffer2 := make([]byte, sizeOfRandString/2)
	_, err = netConn.Read(readBuffer2)
	assert.NoError(t, err)
	received := fmt.Sprintf("%s%s", readBuffer1, readBuffer2)
	assert.Equal(t, string(payload), received)

	sc.Teardown()
}

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func checkReadWrite(sc *ServerClient) error {
	if err := sc.SendSampleLine(); err != nil {
		return err
	}

	if _, err := sc.ReadLine(); err != nil {
		return err
	}
	return nil
}
