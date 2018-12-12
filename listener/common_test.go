// Package listener provides the structs to listen to sockets
package listener_test

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

type netServerModuleMock struct {
	jsonHandler  func([]byte) []byte
	protoHandler func([]byte) []byte
}

func (m *netServerModuleMock) Init(rawConfig yaml.MapSlice) {

}

func (m *netServerModuleMock) HandleJSON(bytes []byte) []byte {
	return m.jsonHandler(bytes)
}

func (m *netServerModuleMock) HandleProto(bytes []byte) []byte {
	return m.protoHandler(bytes)
}

func getEchoModMock(
	t *testing.T,
	expectedJSONResponse []byte,
	expectedProtoResponse []byte,
) *netServerModuleMock {
	return &netServerModuleMock{
		jsonHandler: func(bytes []byte) []byte {
			return bytes
		},
		protoHandler: func(bytes []byte) []byte {
			return bytes
		},
	}
}
