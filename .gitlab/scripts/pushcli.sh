VERSION=$1

gcloud auth configure-docker
echo $GCR_ACCOUNT_KEY > /tmp/$CI_PIPELINE_ID.json
gcloud auth activate-service-account --key-file /tmp/$CI_PIPELINE_ID.json
export GOOGLE_APPLICATION_CREDENTIALS=/tmp/$CI_PIPELINE_ID.json
gcloud config set project red-inspr

gsutil rsync bin gs://inspr-cli/$VERSION