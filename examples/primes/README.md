# Inspr Primes Example

## Description

This folder contains the files necessary to run the primes example in your cluster. This example consists of 3 DApps:

1. generator -> creates random numbers between `[0,MODULAR]` and sends them to the primes_channel_1.
2. filter -> receives values from the channel primes_channel_1 and checks if it's a prime, if true sends the number to the primes_channel_2.
3. printer -> receives values from the cahnnel primes_channel_2 and prints them in the stdout of the pod.

## How to run it:

Make sure the latest Inspr CLI version is installed.

- To get the latest version, run the following command:  
  `go install github.com/inspr/inspr/cmd/inspr`

Before running the inspr command to create the app in your cluster be sure to set the configuration beforehand.

Setting up the cluster ip is essential and it can be done using \
 `insprctl config serverip <url_to_server>`.

Now that everything is set and ready, the simpler way to test it is to run \
`insprctl apply -f general.yaml`.

Alternatively you could install the elements in a modular way using:\
`insprctl apply -k <folder_name>`
