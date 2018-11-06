package listener

import (
	"context"
	"net"
	"net-server/module"
	"os"
)

func ListenUnix(
	ctx context.Context,
	mod module.NetServerModule,
	socketPath string,
	messageType string,
	configPath string,
) error {
	if _, err := os.Stat(socketPath); os.IsExist(err) {
		err := os.Remove(socketPath)

		if err != nil {
			return err
		}
	}

	listener, err := net.Listen("unix", socketPath)

	if err != nil {
		return err
	}

	defer listener.Close()
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	unixListener := NewSocketListener(ctx, listener, mod)

	return unixListener.Listen(messageType, configPath)
}
