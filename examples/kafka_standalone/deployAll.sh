#!/bin/bash


kubectl create namespace kafkasa-apps
kubectl apply -n kafkasa-apps -f ./topic/create/create.yaml
kubectl apply -n kafkasa-apps -f ./reader/reader.yaml
kubectl apply -n kafkasa-apps -f ./writer/writer.yaml