#!/bin/bash
VERSION=$(git describe --always --tags)
mkdir bin

echo $VERSION > version

echo "Building CLI for Linux"
# 386
GOOS=linux GOARCH=386 go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-linux-386-$VERSION ./cmd/insprctl 
# amd64
GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-linux-amd64-$VERSION ./cmd/insprctl
# arm
GOOS=linux GOARCH=arm go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-linux-arm-$VERSION ./cmd/insprctl 
# arm64
GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-linux-arm64-$VERSION ./cmd/insprctl 

echo "Building CLI for Windows"
# amd64
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-windows-amd64-$VERSION.exe ./cmd/insprctl
# 386
GOOS=windows GOARCH=386 go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-windows-386-$VERSION.exe ./cmd/insprctl
# arm
GOOS=windows GOARCH=arm go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-windows-arm-$VERSION.exe ./cmd/insprctl

echo "Building CLI for Darwin"
# amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-darwin-amd64-$VERSION ./cmd/insprctl
# arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.version=$VERSION'" -o ./bin/insprcli-darwin-arm64-$VERSION ./cmd/insprctl
