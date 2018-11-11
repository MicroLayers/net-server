package main

import (
	yaml "gopkg.in/yaml.v2"
)

type EchoModule struct{}

func echo(rawConfig yaml.MapSlice, data []byte) []byte {
	return data
}

func (m *EchoModule) HandleJson(rawConfig yaml.MapSlice, data []byte) []byte {
	return echo(rawConfig, data)
}

func (m *EchoModule) HandleProto(rawConfig yaml.MapSlice, data []byte) []byte {
	return echo(rawConfig, data)
}

var NetServerModule EchoModule
