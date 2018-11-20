// Package listener provides the structs to listen to sockets
package listener_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

type netServerModuleMock struct {
	jsonHandler  func(yaml.MapSlice, []byte) []byte
	protoHandler func(yaml.MapSlice, []byte) []byte
}

func (m *netServerModuleMock) HandleJSON(configMapSlice yaml.MapSlice, bytes []byte) []byte {
	return m.jsonHandler(configMapSlice, bytes)
}

func (m *netServerModuleMock) HandleProto(configMapSlice yaml.MapSlice, bytes []byte) []byte {
	return m.protoHandler(configMapSlice, bytes)
}

func getEchoModMock(
	t *testing.T,
	expectedJSONResponse []byte,
	expectedProtoResponse []byte,
	configMapSlice yaml.MapSlice,
) *netServerModuleMock {
	return &netServerModuleMock{
		jsonHandler: func(confMapSlice yaml.MapSlice, bytes []byte) []byte {
			assert.Equal(t, configMapSlice, confMapSlice)
			return bytes
		},
		protoHandler: func(confMapSlice yaml.MapSlice, bytes []byte) []byte {
			assert.Equal(t, configMapSlice, confMapSlice)
			return bytes
		},
	}
}
