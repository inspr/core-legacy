# Inspr PubSub Example

## Description

This folder contains the files necessary to run a PubSub dApp in your cluster. This examples consists in the creation of a dApp in Insprd that can do the following:

1. Expose a service that receives requests containing messages to be sent
2. Send the received messages to a Discord server
3. Send the receives messages to a Slack server

Each of the previous can be done by making a request to the publisher dApp. A request can be made by using `curl`: 
```zsh
curl -H "Content-Type: application/json" \
  -d '{"message": "<message>"}' \
  inspr.com/publish
```

You can read more about PubSub and this specific example [here](../../docs/pubsub.md).

## How to run it:

Make sure the latest Inspr CLI version is installed. To get the latest version, run the following command:  
```zsh
go install inspr.dev/inspr/cmd/inspr`
```

Before running the Inspr CLI command to create the app in your cluster be sure to set the [configuration](../../docs/readme.md) beforehand.  
Also, this example uses Kafka as the Channels Message Broker (for `controllerprimes` dApp), so it must be first installed in the cluster, and then you must install it on Insprd.  

You can install Kafka in Insprd running the following command from within `/examples/pubsub` folder:
```
insprctl cluster config kafka yamls/kafka.yaml
```

Now that everything is set and ready, run the following to create the "workspace" dApp, which will contain the publisher and subscriber dApps:
```
insprctl apply -f yamls/01.app.yaml
```

Then create its Type, Channel and Nodes:
```
insprctl apply -f yamls/02.ct.yaml
insprctl apply -f yamls/03.ch.yaml
insprctl apply -k yamls/nodes
```

Once that is done, you want to get the name of the Service of the publisher dApp that was created in your Kubernetes Cluster and place it in the `ingress.yaml` file in `/k8s`, in the `service.name` field.
```yaml
...
- path: /publish
    pathType: Prefix
    backend:
        service:
        name: <SERVICE_NAME>
        port:
            number: 80
...
```

Then the ingress for the publisher dApp can be applied in the cluster:
```
k apply -f k8s/ingress.yaml --namespace inspr-apps
```