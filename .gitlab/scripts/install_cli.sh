#!/bin/bash

CLI_VERSION=$(curl -s https://storage.googleapis.com/inspr-cli/version_info)

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

curl $CURL_URL -o inspr
chmod +x inspr 
sudo mv inspr /usr/local/bin
