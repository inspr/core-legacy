#!/bin/bash

CLI_VERSION=$(curl -s https://storage.googleapis.com/inspr-cli/latest-version)
CURL_URL="https://storage.googleapis.com/inspr-cli/${CLI_VERSION}/insprcli"

OS_NAME=$(uname -s)
echo 'Your operating system is '$OS_NAME

case "${OS_NAME}" in
    Linux*)     CURL_URL=$CURL_URL"-linux";;
    Darwin*)    CURL_URL=$CURL_URL"-darwin";;
    CYGWIN* | MINGW* | Windows*)
        echo "For windows system trying to run the bash script, please download the executable from the release page"
        exit 1
    ;;
    *)          
        echo "ERROR identifying the os"
        exit 1
    ;;
esac

ARCH=$(uname -p)
echo 'Your computer architecture is '$ARCH

case "${ARCH}" in
    x86_64*) 
        CURL_URL=$CURL_URL"-amd64"
    ;;
    
    amd64*) 
        CURL_URL=$CURL_URL"-amd64"
    ;;
    
    *86)
        if [[ $OS_NAME == Darwin* ]]; then
            CURL_URL=$CURL_URL"-amd64"
        else
            CURL_URL=$CURL_URL"-386"
        fi
    ;;
    
    arm) 
        if [[ $OS_NAME == Linux* ]]; then
            CURL_URL=$CURL_URL"-arm"
        else
            CURL_URL=$CURL_URL"-arm64"
        fi
    ;;

    arm*) 
        CURL_URL=$CURL_URL"-arm64"
    ;;
    
    aarch64)
        CURL_URL=$CURL_URL"-arm64"
    ;;

    armv8*) 
        CURL_URL=$CURL_URL"-arm64"
    ;;
    
    *)  
        echo "ERROR identifying the architecture"
        exit 2
    ;;
esac

# adding the version to the curl URL
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
