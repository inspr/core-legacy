#!/bin/bash
VERSION=$(cat version)

gcloud auth configure-docker
echo ${{ secret.GCR_ACCOUNT_KEY }} > ${{ github.run_id }}.json
gcloud auth activate-service-account --key-file ${{ github.run_id }}.json
export GOOGLE_APPLICATION_CREDENTIALS=${{ github.run_id }}.json
gcloud config set project red-inspr

gsutil rsync bin gs://inspr-cli/$VERSION

echo $VERSION > latest-version
gsutil cp latest-version gs://inspr-cli/latest-version
