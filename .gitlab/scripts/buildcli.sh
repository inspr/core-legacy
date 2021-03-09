#!/bin/bash
VERSION=$(git describe --always)
mkdir bin

echo $VERSION > version

echo "Building CLI for Linux"
GOOS=linux GOARCH=arm64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-linux-arm64-$VERSION ./cmd/inspr 
GOOS=linux GOARCH=amd64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-linux-amd64-$VERSION ./cmd/inspr

echo "Building CLI for Windows"
GOOS=windows GOARCH=amd64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-windows-amd64-$VERSION ./cmd/inspr
GOOS=windows GOARCH=386 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-windows-386-$VERSION ./cmd/inspr

echo "Building CLI for Mac"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-darwin-amd64-$VERSION ./cmd/inspr
