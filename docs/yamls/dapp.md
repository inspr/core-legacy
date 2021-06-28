# dApps

## Definitions

| Field                | Meaning                                                                                                                                                                                     |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion           | Specify what version of the API to use, for example `v1`                                                                                                                                    |
| kind                 | Specifies which structure the file represents, in this case it would be `dapp`                                                                                                              |
| meta                 | Metadata of dApp                                                                                                                                                                            |
| &rarr;name           | Defines dApp name                                                                                                                                                                           |
| &rarr;reference      | String that contains the url to the location of the DApp definition in inspr's registry                                                                                                     |
| &rarr;annotations    | Definitions that can describe characteristics of the dApp that later on can be used to process/group the dApps in your cluster.                                                             |
| &rarr;parent         | Defines dApp context in relation to the cluster. For example `app1.app2` would mean that this dApp is located on the path `root->app1->app2->app-name`. It is injected by the Inspr daemon. |
| spec                 | specification of dApp                                                                                                                                                                       |
| &rarr;node           | Section describing the dApp Node                                                                                                                                                            |
| &rarr;&rarr;meta     | Metadata of the Node                                                                                                                                                                        |
| name                 | Defines Node name                                                                                                                                                                           |
| reference            | String that contains the url to the location of the Node definition in inspr's registry                                                 |
| annotations          | Definitions that can describe characteristics of the Node that later on can be used to process/group the Nodes in your cluster.                                                             |
| parent               | Defines the Node context in relation to the cluster for example `app1.app2` would mean that this Node is located on the path `root->app1->app2`. It is injected by the Inspr daemon.        |
| &rarr;&rarr; spec    |                                                                                                                                                                                             |
| ports                | array of NodePort structure, represents the connections of a node                                                                                                                           |
| &rarr;port           | Port that is used for receiving a message in the node                                                                                                                                       |
| &rarr;targetPort     | definition on which will be the port where the node's message is sent to                                                                                                                    |
| image                | An URL that serve to point to the location in which the docker image of your application is stored                                                                                          |
| replicas             | Defines the amount of replicas to be created in your cluster                                                                                                                                |
| environment          | Defines the environment variables of your Node                                                                                                                                              |
| &rarr; apps          | Set of dApps that are connected to this dApp, can be specified when creating a new dApp or modified when a dApp is updated.                                                                 |
| &rarr; channels      | Set of Channels that are created in the context of this dApp                                                                                                                                |
| &rarr; types         | Set of Types that are created in the context of this dApp                                                                                                                                   |
| &rarr; boundary      |                                                                                                                                                                                             |
| &rarr; &rarr; input  | List of Channels that are used for the input of this dApp                                                                                                                                   |
| &rarr; &rarr; output | List of Channels that are used as the output of this dApp                                                                                                                                   |
| sidecarPort          | Strucuter that specifies the ports used to talk with the load balancer                                                                                                                      |
| &rarr; lbread        | The port in which the load balancer's message will be read from in the node                                                                                                                 |
| &rarr; lbwrite       | The port in which the node will send messages to communicate with the loadbalancer                                                                                                          |

## YAML example
```yaml
apiVersion: v1
kind: dapp
meta:
  name: generator  
spec:
  node:
    meta:
      name: node-generator
      parent: generator
    spec:
      image: gcr.io/red-inspr/inspr/examples/primes/generator:latest
      replicas: 3
      environment:
        MODULE: 100
  boundary:
    input:
      - primes_ch2
    output:
      - primes_ch1   
```

[back](index.md)