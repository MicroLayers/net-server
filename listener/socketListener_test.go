package listener_test

import (
	"context"
	"encoding/binary"
	"net"
	"net-server/listener"
	"net-server/module"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func setupClientServer(t *testing.T) (net.Listener, net.Conn) {
	tcpListener, err := net.Listen("tcp", "localhost:0")
	assert.NoError(t, err)
	clientConn, err := net.Dial("tcp", tcpListener.Addr().String())
	assert.NoError(t, err)

	return tcpListener, clientConn
}

func TestSocketListener_WillListenForJSONMessages(t *testing.T) {
	expectedResponse := []byte{1, 2, 3, 4, 5}
	configMapSlice := yaml.MapSlice{yaml.MapItem{Key: "var", Value: "val"}}
	modMock := getEchoModMock(t, expectedResponse, []byte{}, configMapSlice)

	tcpListener, clientConn := setupClientServer(t)

	ctx, cancelContext := context.WithCancel(context.Background())
	listener := listener.NewSocketListener(
		ctx,
		tcpListener,
		modMock,
	)

	waitGroup := sync.WaitGroup{}
	var listenerErr error

	waitGroup.Add(1)
	go func() {
		listenerErr = listener.Listen(module.MessageTypeJSON, configMapSlice)

		waitGroup.Done()
	}()

	sizeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(sizeBytes, uint16(len(expectedResponse)))
	clientConn.Write(append(sizeBytes, expectedResponse...))
	response := make([]byte, 64)
	readBytesNum, err := clientConn.Read(response)

	assert.Nil(t, err)
	assert.Equal(t, len(expectedResponse), readBytesNum)

	cancelContext()
	clientConn.Close()
	tcpListener.Close()

	waitGroup.Wait()

	assert.NoError(t, listenerErr, "The listener should exit without errors")
	assert.Equal(
		t,
		expectedResponse,
		response[:readBytesNum],
		"The response should be the result of the listener's JSON handler",
	)
}

func TestSocketListener_WillListenForProtoMessages(t *testing.T) {
	expectedResponse := []byte{42, 92, 73, 54, 7}
	configMapSlice := yaml.MapSlice{yaml.MapItem{Key: "var", Value: "val"}}

	modMock := getEchoModMock(t, []byte{}, expectedResponse, configMapSlice)

	tcpListener, clientConn := setupClientServer(t)

	ctx, cancelContext := context.WithCancel(context.Background())
	listener := listener.NewSocketListener(
		ctx,
		tcpListener,
		modMock,
	)

	waitGroup := sync.WaitGroup{}
	var listenerErr error

	waitGroup.Add(1)
	go func() {
		listenerErr = listener.Listen(module.MessageTypeProto, configMapSlice)

		waitGroup.Done()
	}()

	sizeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(sizeBytes, uint16(len(expectedResponse)))
	clientConn.Write(append(sizeBytes, expectedResponse...))
	response := make([]byte, 64)
	readBytesNum, err := clientConn.Read(response)

	assert.Nil(t, err)
	assert.Equal(t, len(expectedResponse), readBytesNum)

	cancelContext()
	clientConn.Close()
	tcpListener.Close()

	waitGroup.Wait()

	assert.NoError(t, listenerErr, "The listener should exit without errors")
	assert.Equal(
		t,
		expectedResponse,
		response[:readBytesNum],
		"The response should be the result of the listener's Proto handler",
	)
}
