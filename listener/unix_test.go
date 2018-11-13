package listener_test

import (
	"context"
	"encoding/binary"
	"io/ioutil"
	"net"
	"net-server/listener"
	"net-server/module"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestUnixListenerWillListenOnUnixSocket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	expectedResponse := []byte{1, 2, 3, 4, 5}
	configMapSlice := yaml.MapSlice{yaml.MapItem{Key: "var", Value: "val"}}
	messageType := module.MessageTypeJSON
	tmpDir, err := ioutil.TempDir(os.TempDir(), "net-server_")
	assert.NoError(t, err, "The test should successfully create a temporary directory")
	defer os.RemoveAll(tmpDir)
	socketPath := filepath.Join(tmpDir, "socket.sock")

	// Create the socket file, as if a previous instance
	// crashed leaving the socket file
	oldFile, err := os.OpenFile(socketPath, os.O_RDONLY|os.O_CREATE, 0666)
	assert.NoError(t, err, "The test should succesfully create the old socket file")
	oldFile.Close()

	mod := getEchoModMock(t, expectedResponse, expectedResponse, configMapSlice)

	go func() {
		err = listener.ListenUnix(ctx, mod, socketPath, messageType, configMapSlice)
		assert.NoError(t, err, "The listener should terminate in a clean way")
	}()

	// Wait the listener to create the socket file
	time.Sleep(time.Second)
	clientConn, err := net.Dial("unix", socketPath)
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
