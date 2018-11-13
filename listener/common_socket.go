// Package listener provides the structs to listen to sockets
package listener

import (
	"bufio"
	"context"
	"encoding/binary"
	"net"
	"net-server/module"
	"sync"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// SocketListener the socket listener struct
type SocketListener struct {
	Context  context.Context
	Listener net.Listener
	Module   module.NetServerModule
}

// Listen listen for messages of the given type, passing the MapSlice to the module
func (l *SocketListener) Listen(messageType string, configMapSlice yaml.MapSlice) error {
	messageProcessingWaitGroup := sync.WaitGroup{}
	shouldContinue := true
	for shouldContinue {
		conn, err := l.Listener.Accept()
		if err != nil {
			if l.Context.Err() == context.Canceled {
				messageProcessingWaitGroup.Wait()

				return nil
			}

			return err
		}
		defer conn.Close()
		go func() {
			<-l.Context.Done()
			conn.Close()
		}()

		go func(conn net.Conn) {
			reader := bufio.NewReader(conn)
			sizeBuffer := make([]byte, 2)
			messageBuffer := make([]byte, 65535)
			var size uint16

			for {
				// First read the bytes length of the message
				readBytesNum, err := conn.Read(sizeBuffer)
				if err != nil {
					break
				}

				if readBytesNum != 2 {
					log.
						WithField("readBytes", readBytesNum).
						Error("Expected message size of length 2")

					break
				}
				size = binary.LittleEndian.Uint16(sizeBuffer)

				readBytesNum, _ = reader.Read(messageBuffer)
				if uint16(readBytesNum) != size {
					log.
						WithField("readBytes", readBytesNum).
						Errorf("Expected message size of length %d", size)

					shouldContinue = false
					break
				}

				message := messageBuffer[:readBytesNum]

				messageProcessingWaitGroup.Add(1)

				var response []byte
				switch messageType {
				case module.MessageTypeJSON:
					response = l.Module.HandleJson(configMapSlice, message)
				case module.MessageTypeProto:
					response = l.Module.HandleProto(configMapSlice, message)
				}

				if len(response) > 0 {
					_, err = conn.Write(response)

					if err != nil {
						log.
							WithField("error", err).
							Error("Error writing the response")
					}
				}

				messageProcessingWaitGroup.Done()
			}
		}(conn)
	}

	messageProcessingWaitGroup.Wait()

	return nil
}

// NewSocketListener create a SocketListener struct using the given arguments
func NewSocketListener(
	ctx context.Context,
	listener net.Listener,
	mod module.NetServerModule,
) SocketListener {
	return SocketListener{
		Context:  ctx,
		Listener: listener,
		Module:   mod,
	}
}
