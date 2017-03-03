#!/bin/bash
env GOOS=linux GOARCH=amd64 go build -o bin/asnbuilder.linux main.go
env GOOS=windows GOARCH=amd64 go  build -o bin/asnbuilder.exe main.go
env GOOS=darwin GOARCH=amd64 go  build -o bin/asnbuilder.mac main.go