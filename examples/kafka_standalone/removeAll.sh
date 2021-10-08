#!/bin/bash

kubectl apply -n kafkasa -f ./topic/create/create.yaml
sleep 10s
kubectl delete namespace kafkasa