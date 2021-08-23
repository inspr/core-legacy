# Inspr Controller dApp Example

## Description

This folder contains the files necessary to run the controller dApp example in your cluster. This examples consists in the creation of a dApp in Insprd that can do the following:

1. Create a new dApp called `controllerprimes` in the root. This dApp is the same as the one defined in `examples/primes` (generates numbers and prints the ones that are primes).
2. Delete the dApp `controllerprimes`.
3. Update the dApp `controllerprimes` adding a new Annotation to it.

Each of the previous can be done by making a request to the controller dApp. A request can be made by using `curl`:  
```zsh
curl inspr.com/<create | update | delete>
```

## How to run it:

Make sure the latest Inspr CLI version is installed.

- To get the latest version, run the following command:  
  `go install inspr.dev/inspr/cmd/inspr`

Before running the Inspr CLI command to create the app in your cluster be sure to set the [configuration](../../docs/readme.md) beforehand.  
Also, this example uses Kafka as the Channels Message Broker (for `controllerprimes` dApp), so it must be first installed in the cluster, and then you must install it on Insprd.  

You can install Kafka in Insprd running the following command from within `/examples/controller` folder:
`insprctl brokers kafka yamls/kafka.yaml`. 

Now that everything is set and ready, run the following to create the Controller dApp:
```
insprctl apply -f yamls/dapp.yaml
```

Once that is done, you want to get the name of the Service that was created in your Kubernetes Cluster and place it in the `ingress.yaml` file in `/yamls`, in all the `paths`.
```yaml
- path: /<create | update | delete>
    pathType: Prefix
    backend:
        service:
        name: <SERVICE_NAME>
        port:
            number: 8080
```

Then the ingress for the Controller dApp can be applied in the cluster:
```
k apply -f yamls/ingress.yaml --namespace inspr-apps
```

Now everything is set! You can now make the requests to the Controller dApp so it creates/updates/deletes the `controllerprimes` dApp!