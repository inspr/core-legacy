# Sidecar demo  

## Description

This demo creates a dApp table that contains two nodes, a ping and a pong, that exchange messages with each other.

## How to run it  

To run this demo, have a running inspr cluster.
Run `make` and run `insprclt apply` command for the YAML in `ctype/`, then apply the YAMLs in `channels/`, then apply `table.yaml` and finally apply the YAMLs in `/nodes` inside this directory.

