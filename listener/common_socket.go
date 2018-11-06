package listener

import (
	"bufio"
	"context"
	"encoding/binary"
	"net"
	"net-server/module"
	"sync"

	log "github.com/sirupsen/logrus"
)

type SocketListener struct {
	Context  context.Context
	Listener net.Listener
	Module   module.NetServerModule
}

func (l *SocketListener) Listen(messageType string, configPath string) error {
	messageProcessingWaitGroup := sync.WaitGroup{}
	for true {
		conn, err := l.Listener.Accept()
		if err != nil {
			if l.Context.Err() == context.Canceled {
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

			for true {
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

				readBytesNum, err = reader.Read(messageBuffer)
				if uint16(readBytesNum) != size {
					log.
						WithField("readBytes", readBytesNum).
						Errorf("Expected message size of length %d", size)

					break
				}

				message := messageBuffer[:readBytesNum]

				messageProcessingWaitGroup.Add(1)

				var response []byte
				switch messageType {
				case module.MessageTypeJson:
					response = l.Module.HandleJson(configPath, message)
					break
				case module.MessageTypeProto:
					response = l.Module.HandleProto(configPath, message)
					break
				}

				if len(response) > 0 {
					conn.Write(response)
				}

				messageProcessingWaitGroup.Done()
			}
		}(conn)
	}

	messageProcessingWaitGroup.Wait()

	return nil
}

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
