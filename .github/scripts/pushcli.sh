#!/bin/bash
GCR_DETAILS=$1
CI_RUN_ID=$2
VERSION=$(cat version)

gcloud auth configure-docker
echo $GCR_DETAILS > ($CI_RUN_ID).json
gcloud auth activate-service-account --key-file ($CI_RUN_ID).json
export GOOGLE_APPLICATION_CREDENTIALS=($CI_RUN_ID).json
gcloud config set project red-inspr

gsutil rsync bin gs://inspr-cli/$VERSION

echo $VERSION > latest-version
gsutil cp latest-version gs://inspr-cli/latest-version
