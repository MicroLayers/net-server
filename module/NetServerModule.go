package module

const MessageTypeJson = "json"
const MessageTypeProto = "proto"

// NetServerPlugin add features to the main server
// exported objects must be named
type NetServerModule interface {
	HandleJson(configurationPath string, bytes []byte) []byte  // ProtoHandler add support for protobuf messages
	HandleProto(configurationPath string, bytes []byte) []byte // JsonHandler add support for JSON messages
}
