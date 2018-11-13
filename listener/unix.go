package listener

import (
	"context"
	"net"
	"net-server/module"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// ListenUnix listen to a Unix socket,
// managing the messages via the provided module
func ListenUnix(
	ctx context.Context,
	mod module.NetServerModule,
	socketPath string,
	messageType string,
	configMapSlice yaml.MapSlice,
) error {
	// Remove socket file if exists
	// net.Listen will fail in case socket file can't be deleted
	os.Remove(socketPath)

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

	return unixListener.Listen(messageType, configMapSlice)
}
