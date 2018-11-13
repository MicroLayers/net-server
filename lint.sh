#!/usr/bin/env bash

set -e
export PATH=$HOME/go/bin:$PATH

hasExecutable=$(which gometalinter 2>/dev/null)
if [ "$hasExecutable" == "" ]; then
  curl -L https://git.io/vp6lP | sh
  mv ./bin/* "$FOME/go/bin/"
  pushd $HOME
    go get github.com/davecgh/go-spew/spew
    go get github.com/pmezard/go-difflib/difflib
    go get github.com/stretchr/testify/assert
    go get golang.org/x/sys/unix
  popd
fi

gometalinter ./...
