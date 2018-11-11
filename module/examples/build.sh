#!/usr/bin/env bash

BUILD_ARGS="$BUILD_ARGS"
PWD="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

pushd $PWD
	go build $BUILD_ARGS -buildmode=plugin -o echo-module.so echo/EchoModule.go
	go build $BUILD_ARGS -buildmode=plugin -o invalid-module-no-symbol.so invalid/InvalidModuleNoNetServerModuleSymbol.go
	go build $BUILD_ARGS -buildmode=plugin -o invalid-module-wrong-symbol-type.so invalid/InvalidModuleWrongSymbol.go
popd