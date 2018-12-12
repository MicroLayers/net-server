package main

import (
	yaml "gopkg.in/yaml.v2"
)

type EchoModule struct{}

func echo(data []byte) []byte {
	return data
}

func (m *EchoModule) Init(rawConfig yaml.MapSlice) {

}

func (m *EchoModule) HandleJSON(data []byte) []byte {
	return echo(data)
}

func (m *EchoModule) HandleProto(data []byte) []byte {
	return echo(data)
}

var NetServerModule EchoModule
