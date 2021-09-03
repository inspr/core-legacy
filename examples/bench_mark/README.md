# Sidecar demo  

## Description

This demo creates a dApp that contains two nodes, a sender and a receiver, that exchange messages with maximum throughput.

## How to run it  

Make sure the latest Inspr CLI version is installed and Insprd running.

- To get the latest version, run the following command:  
  `go install inspr.dev/inspr/cmd/inspr`

Before running the Inspr CLI command to create the app in your cluster be sure to set the [configuration](../../docs/readme.md) beforehand.  
Also, this example uses Kafka as the Channels Message Broker, so it must be first installed in the cluster, and then you must install it on Insprd.  

You can install Kafka in Insprd running the following command from within `/examples/bench_mark` folder:
`insprctl brokers kafka yamls/kafkaConfig.yaml`.

Now that everything is set and ready, run the following commands from within `/bench_mark` to apply the structures in defined in `/yamls` folder:
```
insprctl apply -k yamls/ctypes
insprctl apply -k yamls/channels
insprctl apply -f yamls/bench.yaml
insprctl apply -k yamls/nodes
```