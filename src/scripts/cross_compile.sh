#!/bin/bash
cd ..
env GOOS=linux GOARCH=amd64 go build -o bin/asnbuilder.linux main.go
