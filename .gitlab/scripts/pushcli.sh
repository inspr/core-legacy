#!/bin/bash
VERSION=$(cat version)

echo $GCR_ACCOUNT_KEY

gcloud auth configure-docker
echo $GCR_ACCOUNT_KEY > $CI_PIPELINE_ID.json
gcloud auth activate-service-account --key-file $CI_PIPELINE_ID.json
export GOOGLE_APPLICATION_CREDENTIALS=$CI_PIPELINE_ID.json
gcloud config set project red-inspr

gsutil rsync bin gs://inspr-cli/$VERSION

echo $VERSION > latest-version
gsutil cp latest-version gs://inspr-cli/latest-version
gsutil cp install_cli.sh gs://inspr-cli/install_cli.sh