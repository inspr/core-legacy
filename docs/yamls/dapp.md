# dApps

## Definitions

| Field                | Meaning                                                                                                                                                                                     |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion           | Specify what version of the API to use, for example `v1`                                                                                                                                    |
| kind                 | Specifies which structure the file represents, in this case it would be `dapp`                                                                                                              |
| meta                 | Metadata of dApp                                                                                                                                                                            |
| &rarr;name           | Defines dApp name                                                                                                                                                                           |
| &rarr;reference      | String that is utilized to define certain tags to the dApp in question, a way for the user to categorize the numerous dApps in the cluster.                                                 |
| &rarr;Annotations    | Definitions that can describe characteristics of the dApp that later on can be used to process/group the dApps in your cluster.                                                             |
| &rarr;parent         | Defines dApp context in relation to the cluster. For example `app1.app2` would mean that this dApp is located on the path `root->app1->app2->app-name`. It is injected by the Inspr daemon. |
| &rarr;sha256         | Tags images with their sha256 digest.                                                                                                                                                       |
| spec                 | specification of dApp                                                                                                                                                                       |
| &rarr;Node           | Section describing the dApp Node                                                                                                                                                            |
| &rarr;&rarr;Meta     | Metadata of the Node                                                                                                                                                                        |
| name                 | Defines Node name                                                                                                                                                                           |
| reference            | String that is utilized to define certain tags to the Node in question, a way for the user to categorize the numerous Nodes in the cluster.                                                 |
| Annotations          | Definitions that can describe characteristics of the Node that later on can be used to process/group the Nodes in your cluster.                                                             |
| parent               | Defines the Node context in relation to the cluster for example `app1.app2` would mean that this Node is located on the path `root->app1->app2`. It is injected by the Inspr daemon.        |
| sha256               | tags images with their sha256 digest.                                                                                                                                                       |
| &rarr;&rarr; Spec    |                                                                                                                                                                                             |
| Image                | An URL that serve to point to the location in which the docker image of your application is stored                                                                                          |
| Replicas             | Defines the amount of replicas to be created in your cluster                                                                                                                                |
| Environment          | Defines the environment variables of your Node                                                                                                                                              |
| &rarr; Apps          | Set of dApps that are connected to this dApp, can be specified when creating a new dApp or modified when a dApp is updated.                                                                 |
| &rarr; Channels      | Set of Channels that are created in the context of this dApp                                                                                                                                |
| &rarr; Types  | Set of Types that are created in the context of this dApp                                                                                                                           |
| &rarr; Boundary      |                                                                                                                                                                                             |
| &rarr; &rarr; Input  | List of Channels that are used for the input of this dApp                                                                                                                                   |
| &rarr; &rarr; Output | List of Channels that are used as the output of this dApp                                                                                                                                   |

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