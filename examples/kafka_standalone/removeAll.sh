#!/bin/bash

kubectl apply -n kafkasa-apps -f ./topic/delete/delete.yaml
sleep 5s
kubectl delete namespace kafkasa-apps