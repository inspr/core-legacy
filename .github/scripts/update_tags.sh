#!/bin/bash
CHART_URL=$1
CHART_GS_URI=$2
echo "Updating helm chart tags"

OUTPUT=$(cat tags.out | jq ".builds[]" --compact-output)
for k in ${OUTPUT[@]}; do
    IMAGE=`printf %s $(echo $k | jq .imageName)`
    IMAGE=$(eval echo $IMAGE)
    echo "Setting $IMAGE tag to : $TAG"
    TAG=`printf %s $(echo $k | jq .tag)`
    TAG=$(eval echo $TAG)
    sed -i s_${IMAGE}_${TAG}_g ./build/helm/values.yaml
done


APP_VERSION=$(git describe --always --tags)
echo "Updating chart app version to $APP_VERSION"
sed -i 's/appVersion: .*/appVersion: '"$APP_VERSION"'/' build/helm/Chart.yaml

echo "Updating chart version to $APP_VERSION"
sed -i 's/version: .*/version: '"$APP_VERSION"'/' build/helm/Chart.yaml


echo "Updating helm dependencies..."
helm dependency update build/helm

echo "Generating new helm charts"
helm package build/helm -d charts
echo "Helm charts generated"


echo "Updating inspr helm repositories"
gsutil -h "Cache-Control:no-cache,max-age=0" rsync  $CHART_GS_URI charts

helm repo index charts --url $CHART_URL
gsutil -h "Cache-Control:no-cache,max-age=0" rsync  charts $CHART_GS_URI
echo "Helm chart repo updated"
