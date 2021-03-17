#!/bin/bash

CLI_VERSION=$(curl -s https://storage.googleapis.com/inspr-cli/latest-version)
CURL_URL="https://storage.googleapis.com/inspr-cli/${CLI_VERSION}/insprcli"

OS_NAME=$(uname -o)
if [[ $OS_NAME == 'GNU/Linux' ]]; then # windows
    CURL_URL=$CURL_URL"-linux"
elif [[ $OS_NAME == 'Windows' ]]; then # linux
    CURL_URL=$CURL_URL"-windows"
elif  [[ $OS_NAME == 'Darwin' ]]; then  # darwin
    CURL_URL=$CURL_URL"-darwin"
fi

ARCH=$(uname -i)
if [[ $ARCH == x86_64* ]]; then # x64
    CURL_URL=$CURL_URL"-amd64"
elif [[ $ARCH == i*86 ]]; then # x32
    CURL_URL=$CURL_URL"-386"
elif  [[ $ARCH == arm* ]] || [[ $ARCH = aarch64 ]]; then  # arm
    CURL_URL=$CURL_URL"-arm64"
fi

CURL_URL=$CURL_URL"-"$CLI_VERSION

echo $CURL_URL
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