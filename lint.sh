#!/usr/bin/env bash

set -e
set -o xtrace

export BINDIR=$HOME/go/bin
export PATH="${BINDIR}:$PATH"
echo "Installing gometalinter..."
curl -L https://git.io/vp6lP | sh

pushd $HOME/go
  echo "Installing gometalinter deps..."
  go get -u github.com/davecgh/go-spew/spew
  go get -u github.com/pmezard/go-difflib/difflib
  go get -u github.com/stretchr/testify/assert
  go get -u golang.org/x/sys/unix
popd

echo "Executing gometalinter..."
gometalinter --fast --vendor ./...
