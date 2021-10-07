#!/bin/bash


kubectl create namespace kafkasa
kubectl apply -n kafkasa -f ./topic/create/create.yaml
kubectl apply -n kafkasa -f ./reader/reader.yaml
kubectl apply -n kafkasa -f ./writer/writer.yaml