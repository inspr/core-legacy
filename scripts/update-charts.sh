#!/usr/bin/env bash
echo
"
 This script generates a new .tgz file for the charts in
 build/charts and publish them in the inspr gcloud.

 Make sure you have helm and gsutil installed.

"

CHARTS=(/inspr-stack /inspr-stack/subcharts/insprd /inspr-stack/subcharts/uidp)

echo "Creating new .tgz file..."
for chart in ${CHARTS[@]}; do
    helm dependency update ../build$chart
    helm package ../build$chart -d ../build/charts
done

echo "Creating updated index.yaml file..."
helm repo index ../build/charts --url https://inspr-charts.storage.googleapis.com

echo "Synchronizing the update chart(s) to google cloud..."
./sync-repo.sh ../build/charts inspr-charts

echo "Removing generated .tgz chart files"
rm ../build/charts/*.tgz