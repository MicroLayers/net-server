package module

import (
	"fmt"
	"plugin"

	yaml "gopkg.in/yaml.v2"
)

const MessageTypeJson = "json"
const MessageTypeProto = "proto"

// NetServerPlugin add features to the main server
// exported objects must be named
type NetServerModule interface {
	HandleJson(yaml.MapSlice, []byte) []byte  // ProtoHandler add support for protobuf messages
	HandleProto(yaml.MapSlice, []byte) []byte // JsonHandler add support for JSON messages
}

func LoadModule(path string) (NetServerModule, error) {
	plug, err := plugin.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Impossible to open file %s", path)
	}

	symPlugin, err := plug.Lookup("NetServerModule")

	if err != nil {
		return nil, fmt.Errorf(
			"Impossible to find the symbol 'NetServerModule' in %s",
			path,
		)
	}

	mod, ok := symPlugin.(NetServerModule)

	if !ok {
		return nil, fmt.Errorf(
			"Module %s does not implement the required interface",
			path,
		)
	}

	return mod, nil
}
