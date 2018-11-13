#!/usr/bin/env bash

set -e
export PATH=$HOME/go/bin:$PATH

go get -u golang.org/x/lint/golint
configuration=$(golint configuration)
module=$(golint module)
listener=$(golint listener)

echo $configuration
echo $module
echo $listener

if [ "$configuration" != "" ] || [ "$module" != "" ] || [ "$listener" != "" ]; then
  exit 255
fi

exit 0
