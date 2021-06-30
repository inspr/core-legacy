# Sidecar demo  

## Description

This demo creates a dApp table that contains two nodes, a ping and a pong, that exchange messages with each other.

## How to run it  

Make sure the latest Inspr CLI version is installed and Insprd running.

- To get the latest version, run the following command:  
  `go install inspr.dev/inspr/cmd/inspr`

Before running the Inspr CLI command to create the app in your cluster be sure to set the [configuration](../../docs/readme.md) beforehand.  
Also, this example uses Kafka as the Channels Message Broker, so it must be first installed in the cluster, and then you must install it on Insprd.  

You can install Kafka in Insprd running the following command from within `/examples/pingpong_demo` folder:
`insprctl cluster config kafka yamls/kafka.yaml`.

Now that everything is set and ready, run the following commands from within `/pingpong_demo` to apply the structures in defined in `/yamls` folder:
```
insprctl apply -k yamls/ctypes
insprctl apply -k yamls/channels
insprctl apply -f yamls/table.yaml
insprctl apply -k yamls/nodes
```