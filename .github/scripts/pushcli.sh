#!/bin/bash
VERSION=$(cat version)

gsutil rsync bin gs://inspr-cli/$VERSION

echo $VERSION > latest-version
gsutil cp latest-version gs://inspr-cli/latest-version
