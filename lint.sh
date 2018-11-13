#!/usr/bin/env bash

set -e
export PATH="$(pwd)/bin:$PATH"

curl -L https://git.io/vp6lP | sh
pushd $HOME
  go get -u github.com/davecgh/go-spew/spew
  go get -u github.com/pmezard/go-difflib/difflib
  go get -u github.com/stretchr/testify/assert
  go get -u golang.org/x/sys/unix
popd

gometalinter ./...
