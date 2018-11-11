#!/usr/bin/env bash

export BUILD_ARGS="-race"

# First compile the echo module (required for module tests)
PWD="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
pushd $PWD/module/examples
	./build.sh
popd

mkdir -p build
go test $BUILD_ARGS -cover ./... -covermode=atomic -coverprofile=build/coverage.txt
go tool cover -html=build/coverage.txt -o build/coverage.html
