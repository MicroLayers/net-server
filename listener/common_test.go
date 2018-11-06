package listener_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type netServerModuleMock struct {
	jsonHandler  func(string, []byte) []byte
	protoHandler func(string, []byte) []byte
}

func (m *netServerModuleMock) HandleJson(configurationPath string, bytes []byte) []byte {
	return m.jsonHandler(configurationPath, bytes)
}

func (m *netServerModuleMock) HandleProto(configurationPath string, bytes []byte) []byte {
	return m.protoHandler(configurationPath, bytes)
}

func getEchoModMock(
	t *testing.T,
	expectedJSONResponse []byte,
	expectedProtoResponse []byte,
	configPath string,
) *netServerModuleMock {
	return &netServerModuleMock{
		jsonHandler: func(confPath string, bytes []byte) []byte {
			assert.Equal(t, configPath, confPath)
			return bytes
		},
		protoHandler: func(confPath string, bytes []byte) []byte {
			assert.Equal(t, configPath, confPath)
			return bytes
		},
	}
}
