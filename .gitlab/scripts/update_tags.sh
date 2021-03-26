#!/bin/bash
CHART_URL=$1
CHART_GS_URI=$2
echo "Updating helm chart tags"
INSPRD_TAG=$(printf '%s' $(grep -Eo "gcr\.io\/red\-inspr\/insprd\:[^@]+\@sha256\:[[:alnum:]]+" tags.out) | sed -e 's/[\/&]/\\&/g')
KAFKA_SIDECAR_TAG=$(printf '%s' $(grep -Eo "gcr\.io\/red\-inspr\/inspr\/sidecar\/kafka\:[^@]+\@sha256\:[[:alnum:]]+" tags.out) | sed -e 's/[\/&]/\\&/g')

echo
echo "Tags for updating:"
echo "$INSPRD_TAG"
echo "$KAFKA_SIDECAR_TAG"

echo "Updating tags"
sed -i 's/gcr\.io\/red\-inspr\/insprd/'"$INSPRD_TAG"'/'  build/helm/values.yaml
sed -i 's/gcr\.io\/red\-inspr\/inspr\/sidecar\/kafka/'"$KAFKA_SIDECAR_TAG"'/' build/helm/values.yaml

APP_VERSION=$(git describe --always --tags)
echo "Updating chart app version to $APP_VERSION"
sed -i 's/appVersion: .*/appVersion: '"$APP_VERSION"'/' build/helm/Chart.yaml


VERSION=$(grep -Po "version: \K.*" build/helm/Chart.yaml)
NEW_VERSION=$(VERSION=$VERSION KIND=patch python .gitlab/scripts/upgrade_version.py)

echo "Updating chart version to $NEW_VERSION"
sed -i 's/version: .*/version: '"$NEW_VERSION"'/' build/helm/Chart.yaml


echo "Updating helm dependencies..."
helm dependency  update build/helm

echo "Generating new helm charts"
helm package build/helm -d charts
echo "Helm charts generated"


echo "Updating inspr helm repositories"
gsutil -h "Cache-Control:no-cache,max-age=0" rsync  $CHART_GS_URI charts

helm repo index charts --url $CHART_URL
gsutil -h "Cache-Control:no-cache,max-age=0" rsync  charts $CHART_GS_URI
echo "Helm chart repo updated"
