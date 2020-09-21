#!/usr/bin/env bash
version=$(($(cat version)+1)); printf ${version} > version
export GOOS=linux GOARCH=amd64
gofmt -w .
go build -ldflags '-X main.BuildID='${version} -o ./bin/proxytcp main.go
