
#!/usr/bin/env bash

mkdir -p build
go test -race -cover ./... -covermode=atomic -coverprofile=build/coverage.txt
go tool cover -html=build/coverage.txt -o build/coverage.html
