package listener_test

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net-server/listener"
	"net-server/module"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func getFreeTCPPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestTCPListenerWillListenOnTCPPort(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	expectedResponse := []byte{1, 2, 3, 4, 5}
	configMapSlice := yaml.MapSlice{yaml.MapItem{Key: "var", Value: "val"}}
	messageType := module.MessageTypeJson
	freePort, err := getFreeTCPPort()
	assert.NoError(t, err)
	listeningAddress := fmt.Sprintf("localhost:%d", freePort)

	mod := getEchoModMock(t, expectedResponse, expectedResponse, configMapSlice)

	go func() {
		err = listener.ListenTCP(ctx, mod, listeningAddress, messageType, configMapSlice)
		assert.NoError(t, err, "The listener should terminate in a clean way")
	}()

	clientConn, err := net.Dial("tcp", listeningAddress)
	assert.NoError(t, err)

	sizeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(sizeBytes, uint16(len(expectedResponse)))
	sentMessage := append(sizeBytes, expectedResponse...)
	numBytesSent, err := clientConn.Write(sentMessage)
	assert.Equal(t, len(sentMessage), numBytesSent, "The client should send exactly 3 bytes of data")
	assert.NoError(t, err)

	buffer := make([]byte, 10)
	numBytesRead, err := clientConn.Read(buffer)
	assert.Equal(t, len(expectedResponse), numBytesRead)
	assert.NoError(t, err, "The server should send the expected response size")

	message := buffer[:numBytesRead]
	assert.Equal(t, expectedResponse, message, "The server should send the expected response back")

	clientConn.Close()
	cancel()
}
