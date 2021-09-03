#!/bin/bash

insprctl brokers kafka yamls/kafkaConfig.yaml
sleep 2

insprctl apply -k yamls/ctypes
sleep 2

insprctl apply -k yamls/channels
sleep 2

insprctl apply -f yamls/table.yaml
sleep 2

insprctl apply -k yamls/nodes
sleep 2