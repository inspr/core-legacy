#!/bin/bash

CLI_VERSION=$(curl -s https://storage.googleapis.com/inspr-cli/latest-version)
CURL_URL="https://storage.googleapis.com/inspr-cli/${CLI_VERSION}/insprcli"

OS_NAME=$(uname -s)
echo 'Your operating system is '$OS_NAME

case "${OS_NAME}" in
    Linux*)     CURL_URL=$CURL_URL"-linux";;
    Darwin*)    CURL_URL=$CURL_URL"-darwin";;
    CYGWIN*)    CURL_URL=$CURL_URL"-windows";;
    MINGW*)     CURL_URL=$CURL_URL"-windows";;
    Windows*)   CURL_URL=$CURL_URL"-windows";;
    *)          echo "ERROR identifing the os"
esac


ARCH=$(uname -i)
echo 'Your computer architecture is '$ARCH

case "${ARCH}" in
    x86_64*) CURL_URL=$CURL_URL"-amd64";;
    i*86) CURL_URL=$CURL_URL"-386";;
    arm*) CURL_URL=$CURL_URL"-arm64";;
    aarch64) CURL_URL=$CURL_URL"-arm64";;
    *)  echo "ERROR identifing the architecture"
esac

CURL_URL=$CURL_URL"-"$CLI_VERSION
echo 'Downloading the inspr cli binary'
curl $CURL_URL -o /tmp/inspr

ENCODING=utf-8
if iconv --from-code="$ENCODING" --to-code="$ENCODING" /tmp/inspr > /dev/null 2>&1; then
    echo "error, coudln't find the binary, in the url used"
    echo $CURL_URL
else    
    chmod +x /tmp/inspr 
    echo 'Moving binary into /usr/local/bin, need sudo permission'
    sudo mv /tmp/inspr /usr/local/bin
    echo 'Files moved to to /usr/local/bin'
fi