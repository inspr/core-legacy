# [BETA] Routes

Inspr allows communication via routes, that is, via http calls between nodes. In a simple and quick way, it is possible to register endpoints and send requests from one dApp to another dApp. This feature is still in the beta phase, but the main features are already up and running. This step of development of the routes feature includes the essential ---- of the routes structure. In this beta, you can register handler functions in pre-determined paths, that should be intantiate when you install the dapp on Inspr. When you're writing the dApp yaml, you can include the endpoints you need for each node you want. An declared endpoint will be available for all the other nodes inside the same parent dApp. An endpoint to be available doesn't mean it exists in the node, so you still need to create it. For that, Inspr have a route client library, that handle paths and send request in a simple way. Lets walkthrough an example to clarify to make it clear.

## Declaring the endpoints in the yaml

So the first thing you need to do is decide what endpoints your apps will have. In this example, we will have two nodes dapps that will comunicate between them. Those nodes will have a parent dapp, called `router`. Inside of it we declare both the `api` and the `client` dapps. Those dapps (the `api` and the `client`) are node dapps, therefore it can't have any children. The idea here is that the `api` will register some endpoints, while the `client` will automatically send requests to those endpoints, and then receiving a response from the `api`. For this example, the `api` will declare three endpoints:

- `add`, that adds two numbers received in the request body.
- `sub`, that substract two numbers received in the request body.
- `mul`, that multiply two numbers received in the request body (as you already guess it).

As the client will only send requests, no endpoint will be declared for it. Simple enough. The final yaml will look like this:

```yaml
apiVersion: v1
kind: dapp

meta:
  name: router
spec:
  apps:
    client:
      meta:
        name: client
        parent: router
      spec:
        node:
          spec:
            replicas: 1
            image: <client_app_image>
    api:
      meta:
        name: api
        parent: router
      spec:
        node:
          spec:
            replicas: 1
            image: <api_app_image>
            endpoints:
              - add
              - sub
              - mul

```
Note that the images fields of both the `client` and the `api` need to be filled, as is the code that will run once the dApp is up in the cluster. In the next section, we'll write this code using Inspr's client library, that can register/listen specific endpoints and send requests to another node.

## Using Inspr's Route Client



