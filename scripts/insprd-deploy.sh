#!/usr/bin/env bash

echo "Adding repository..."
helm repo add inspr https://inspr-charts.storage.googleapis.com/

echo "Update repository..."
helm repo update

echo "Deploying to cluster..."
helm upgrade --install insprd inspr/insprd