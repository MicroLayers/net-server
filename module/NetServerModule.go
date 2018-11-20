package module

import (
	"fmt"
	"plugin"

	yaml "gopkg.in/yaml.v2"
)

// MessageTypeJSON message type JSON
const MessageTypeJSON = "json"

// MessageTypeProto message type protobuf
const MessageTypeProto = "proto"

// NetServerModule add features to the main server
// exported objects must be named
type NetServerModule interface {
	HandleJSON(yaml.MapSlice, []byte) []byte  // HandleJSON add support for JSON messages
	HandleProto(yaml.MapSlice, []byte) []byte // HandleProto add support for protobuf messages
}

// LoadModule load the NetServerModule from the given path
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
