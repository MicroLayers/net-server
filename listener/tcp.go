package listener

import (
	"context"
	"net"
	"net-server/module"

	yaml "gopkg.in/yaml.v2"
)

// ListenTCP listen to a TCP port,
// managing the messages via the provided module
func ListenTCP(
	ctx context.Context,
	mod module.NetServerModule,
	listeningAddress string,
	messageType string,
	configMapSlice yaml.MapSlice,
) error {
	listener, err := net.Listen("tcp", listeningAddress)

	if err != nil {
		return err
	}

	defer listener.Close()
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	tcpListener := NewSocketListener(ctx, listener, mod)

	return tcpListener.Listen(messageType, configMapSlice)
}
