#!/bin/bash
GITHUB_SHA=$1
VERSION=$(cat version)
BASE_URL="https://console.cloud.google.com/storage/browser/inspr-cli/bin"

release-cli create --name "Release $VERSION" \
        --tag-name $VERSION \
        --ref $GITHUB_SHA \
        --assets-link "{\"name\":\"insprcli-linux-arm64\",
                        \"url\":\"$BASE_URL/insprcli-linux-arm64-$VERSION\"}" \
        --assets-link "{\"name\":\"insprcli-linux-amd64\",
                        \"url\":\"$BASE_URL/bin/insprcli-linux-amd64-$VERSION\"}" \
        --assets-link "{\"name\":\"insprcli-windows-amd64\",
                        \"url\":\"$BASE_URL/bin/insprcli-windows-amd64-$VERSION\"}" \
        --assets-link "{\"name\":\"insprcli-windows-386\",
                        \"url\":\"$BASE_URL/bin/insprcli-windows-386-$VERSION\"}" \
        --assets-link "{\"name\":\"insprcli-darwin-amd64\",
                        \"url\":\"$BASE_URL/bin/insprcli-darwin-amd64-$VERSION\"}"