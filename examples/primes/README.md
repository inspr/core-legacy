# Inspr Primes Example

## Description

This folder contains the files necessary to run the primes example in your cluster. This example consists of 3 dApps:

1. generator -> creates random numbers between `[0,MODULAR]` and sends them to the Channel primesch1.
2. filter -> receives values from the Channel primesch1 and checks if it's a prime, if true sends the number to the Channel primesch2.
3. printer -> receives values from the Channel primesch2 and prints them in the stdout of the pod.

## How to run it:

Make sure the latest Inspr CLI version is installed.

- To get the latest version, run the following command:  
  `go install inspr.dev/inspr/cmd/inspr`

Before running the Inspr CLI command to create the app in your cluster be sure to set the [configuration](../../docs/readme.md) beforehand.  
Also, this example uses Kafka as the Channels Message Broker, so it must be first installed in the cluster, and then you must install it on Insprd.  

You can install Kafka in Insprd running the following command from within `/examples/primes` folder:
`insprctl cluster config kafka yamls/kafka.yaml`.  

Now that everything is set and ready, the simpler way to test it is to run:
```
insprctl apply -f yamls/general.yaml
```
This will create the dApp `primesexample` that contains all the needed structures in it (Channels, Types, Aliases and other dApps).

Alternatively you could create the structures in a modular way using the following commands (in this order):
```
insprctl apply -k yamls/types
insprctl apply -k yamls/channels
insprctl apply -f yamls/basedapp.yaml
insprctl apply -k yamls/nodes
```
