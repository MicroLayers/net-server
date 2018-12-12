package listener

import (
	"context"
	"net"
	"net-server/module"
)

// ListenTCP listen to a TCP port,
// managing the messages via the provided module
func ListenTCP(
	ctx context.Context,
	mod module.NetServerModule,
	listeningAddress string,
	messageType string,
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

	return tcpListener.Listen(messageType)
}
