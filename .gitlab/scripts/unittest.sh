#!/bin/bash
go test -v -coverprofile=profile.cov.tmp -coverpkg=./... ./...
cat profile.cov.tmp | grep -v 'fake\|mock\|examples\|main.go' > profile.cov
go tool cover -func profile.cov