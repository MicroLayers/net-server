package listener

import (
	"context"
	"net"
	"net-server/module"
)

func ListenTCP(
	ctx context.Context,
	mod module.NetServerModule,
	listeningAddress string,
	messageType string,
	configPath string,
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

	return tcpListener.Listen(messageType, configPath)
}
