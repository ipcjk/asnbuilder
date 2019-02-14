#!/bin/bash
mkdir release
env GOOS=linux GOARCH=amd64 go build -o release/asnbuilder.linux main.go
env GOOS=windows GOARCH=amd64 go  build -o release/asnbuilder.exe main.go
env GOOS=darwin GOARCH=amd64 go  build -o release/asnbuilder.mac main.go