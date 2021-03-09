#!/bin/bash
VERSION=$1

echo "This echo message is working"

echo $GCR_ACCOUNT_KEY

gcloud auth configure-docker
echo $GCR_ACCOUNT_KEY > $CI_PIPELINE_ID.json
gcloud auth activate-service-account --key-file $CI_PIPELINE_ID.json
export GOOGLE_APPLICATION_CREDENTIALS=$CI_PIPELINE_ID.json
gcloud config set project red-inspr

gsutil rsync bin gs://inspr-cli/$VERSION