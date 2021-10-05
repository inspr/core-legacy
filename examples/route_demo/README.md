# Inspr Routes Example

## Description

This folder contains the files necessary to run the routes example in your cluster. This example consists of 2 dApps:

1. client -> sends requests to the api's route.
2. api -> receives request from the client and performs a simple operation dictatade by the endpoint reached (*add*, *sub* or *mul*).

## How to run it:

Make sure the latest Inspr CLI version is installed.

- To get the latest version, run the following command:  
  `go install inspr.dev/inspr/cmd/inspr`

Before running the Inspr CLI command to create the app in your cluster be sure to set the [configuration](../../docs/readme.md) beforehand.  


This example doesn't use brokers, therefore everything is set and ready and the simpler way to test it is to run:
```
insprctl apply -f yamls/route-demo.yaml
```
This will create the dApp `router` that contains all the needed structures in it (other dApps and Routes).
