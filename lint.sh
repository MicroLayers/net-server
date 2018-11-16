#!/usr/bin/env bash

set -e
set -o xtrace

# create a GOPATH for linters
mkdir -p $HOME/lint
# create a GOPATH for usual source
mkdir -p $HOME/go/{src,bin,pkg}

# link them together, so that they install and build to the same directory
# this also means you can cache `$HOME/go/pkg` to avoid redownloading dependencies
ln -s $HOME/go/pkg $HOME/lint/pkg
ln -s $HOME/go/bin $HOME/lint/bin

# installing linters
export GOBIN=$HOME/go/bin
export GOPATH=$HOME/lint
export GO111MODULE=off
curl -L https://git.io/vp6lP | bash -s -- -b $GOBIN
export GO111MODULE=on

# getting your project linted
export GOPATH=$HOME/go

# download all dependencies and put them into ./vendor
go mod vendor

# flatten vendor to src
cp -r ./vendor/* $HOME/go/src/
rm -rf ./vendor

# all done for running everything with modules off
GO111MODULE=off gometalinter ./...
