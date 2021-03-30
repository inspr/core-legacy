#!/bin/bash
VERSION=$(git describe --always --tags)
mkdir bin

echo $VERSION > version

echo "Building CLI for Linux"
# 386
GOOS=linux GOARCH=386 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-linux-386-$VERSION ./cmd/inspr 
# amd64
GOOS=linux GOARCH=amd64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-linux-amd64-$VERSION ./cmd/inspr
# arm
GOOS=linux GOARCH=arm go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-linux-arm-$VERSION ./cmd/inspr 
# arm64
GOOS=linux GOARCH=arm64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-linux-arm64-$VERSION ./cmd/inspr 

echo "Building CLI for Windows"
# amd64
GOOS=windows GOARCH=amd64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-windows-amd64-$VERSION ./cmd/inspr
# 386
GOOS=windows GOARCH=386 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-windows-386-$VERSION ./cmd/inspr
# arm
GOOS=windows GOARCH=arm go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-windows-arm-$VERSION ./cmd/inspr

echo "Building CLI for Darwin"
# amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-darwin-amd64-$VERSION ./cmd/inspr
# arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X cmd.inspr.version=$VERSION" -o ./bin/insprcli-darwin-arm64-$VERSION ./cmd/inspr
