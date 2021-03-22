# YAML Documentation

The creation of YAML files allows the proper usage of the Inspr CLI.

In this way, when the file is written in one of the formats described below it can be processed and applied to the cluster through the CLI commands `inspr apply -f <file_path>` or `inspr apply -k <files_dir>`.

## dApps

### Definitions

| Field                | Meaning                                                                                                                                                                                     |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion           | specify what version of the API to use, for example `v1`                                                                                                                                    |
| kind                 | specifies which structure the file represents, in this case it would be `dapp`                                                                                                              |
| meta                 | metadata of dApp                                                                                                                                                                            |
| &rarr;name           | defines dApp name                                                                                                                                                                           |
| &rarr;reference      | string that is utilized to defined certain tags to the dApp in question, a way for the user to categorize the numerous dApps in the cluster.                                                |
| &rarr;Annotations    | definitions that can describe characteristics of the dApp that later on can be used to process/group the dApps in your cluster.                                                             |
| &rarr;parent         | defines dApp context in relation to the cluster. For example `app1.app2` would mean that this dApp is located on the path `root->app1->app2->app-name`. It is injected by the inspr daemon. |
| &rarr;sha256         | tags images with their sha256 digest.                                                                                                                                                       |
| spec                 | specification of dApp                                                                                                                                                                       |
| &rarr;Node           | Section describing the dApp node                                                                                                                                                            |
| &rarr;&rarr;Meta     | metadata of the Node                                                                                                                                                                        |
| name                 | defines node name                                                                                                                                                                           |
| reference            | string that is utilized to defined certain tags to the dApp in question, a way for the user to categorize the numerous dApps in the cluster.                                                |
| Annotations          | definitions that can describe characteristics of the node that later on can be used to process/group the nodes in your cluster.                                                             |
| parent               | defines the node context in relation to the cluster for example `app1.app2` would mean that this node is located on the path `root->app1->app2`. It is injected by the inspr daemon.        |
| sha256               | tags images with their sha256 digest.                                                                                                                                                       |
| &rarr;&rarr; Spec    |                                                                                                                                                                                             |
| Image                | an URL that serve to point to the location in which the docker image of your application is stored                                                                                          |
| Replicas             | defines the amount of replicas to be created in your cluster                                                                                                                                |
| Environment          | defines the environment variables of your pods                                                                                                                                              |
| &rarr; Apps          | set of dApps that are connected to this dApp, can be specified when creating a new dApp or modified when a dApp is updated by the inspr daemon                                              |
| &rarr; Channels      | set of Channels that are created in the context of this dApp                                                                                                                                |
| &rarr; ChannelTypes  | set of Channel Types that are created in the context of this dApp                                                                                                                           |
| &rarr; Boundary      |                                                                                                                                                                                             |
| &rarr; &rarr; Input  | List of Channels that are used for the input of this dApp                                                                                                                                   |
| &rarr; &rarr; Output | List of Channels that are used for the output of this dApp                                                                                                                                  |

### YAML example
```yaml
apiVersion: v1
kind: dapp
meta:
  name: "generator"
  reference: ""
  parent: ""
spec:
  node:
    meta:
      name: "node-generator"
      parent: "generator"
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

## Channels 

### Definitions

| Field             | Meaning                                                                                                                                                                                                               |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion        | specify what version of the API to use, for example `"v1"`                                                                                                                                                            |
| kind              | specifies which structure the file represents, in this case it would be `channel`                                                                                                                                     |
| meta              | metadata of Channel                                                                                                                                                                                                   |
| &rarr; name       | defines the Channel name                                                                                                                                                                                              |
| &rarr; reference  | URL reference to the Channel definition in the inspr repository, there are already well defined Channel that can be used instead of defining your own.                                                                |
| &rarr;Annotations | definitions that can describe characteristics of the Channel that later on can be used to process/group the Channels in your cluster.                                                                                 |
| &rarr; parent     | it is injected by the inspr daemon, defines the Channel context in the cluster through the path of the dApp in which the Channel is stored, for example: `app1.app2` means that the Channel is defined in the `app2`. |
| &rarr; sha256     | tags images with their sha256 digest.                                                                                                                                                                                 |
| spec              |                                                                                                                                                                                                                       |
| &rarr; type       | defines the type of the Channel, this field is a string that contains the name of any of the Channel Types defined in your cluster                                                                                    |
| connectedapps     | List of dApp names that are using this Channel, this is injected by the inspr daemon                                                                                                                                  |
| connectedaliasses | A simple list of the aliases that are being used for to reference this channel. Injected by the inspr daemon.                                                                                                         |

### YAML example
```yaml
apiVersion: "v1"
kind: "channel"
meta:
  name: "primes_ch1"
  reference: ""
  Annotations: 
    kafka.partition.number: "1"
    kafka.replication.factor: "1"
  parent: ""  
spec:
  type: "primes_ct1"
```


## Channel Types 

### Definitions

| Field             | Meaning                                                                                                                                                                         |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion        | specify what version of the API to use, for example `"v1"`                                                                                                                      |
| kind              | specifies which structure the file represents, in this case it would be `channeltype`                                                                                           |
| meta              | metadata of Channel_Type                                                                                                                                                        |
| &rarr;name        | channel_type_name                                                                                                                                                               |
| &rarr;reference   | string that is utilized to defined certain tags to the dApp in question, a way for the user to categorize the numerous dApps in the cluster.                                    |
| &rarr;Annotations | definitions that can describe characteristics of the Channel Type that later on can be used to process/group the Channel Types in your cluster.                                 |
| &rarr;parent      | It is injected by the inspr daemon and it's string composed of it's location's path, for example `app1.app2` means that the Channel type belongs to the `app2` in your cluster. |
| &rarr;sha256      | tags images with their sha256 digest.                                                                                                                                           |
| schema            | defines the data structure that goes through this Channel Type, example:  `'{"type":"int"}'`                                                                                    |
| connectedchannels | Is a list of Channels names that are created using this specific type, this is injected through the `inspr CLI`/ `inspr daemon`                                                 |


### YAML example
```yaml
apiVersion: "v1"
kind: "channeltype"
meta:
  name: "primes_ct1"
  reference: ""
  parent: ""  
schema: '{"type":"int"}'
```

## General file

### Definition

The so called general file, or composed file, is nothing more than a YAML that congregates two or more definitions of one of the elements described above. 

For example an dApp that has a collection of other dApps plus some definitions of channel Types and channels.

### YAML example

```yaml
apiVersion: v1
kind: dapp
meta:
  name: "basic-example"
  reference: ""
  parent: ""
  Annotations: 
    kafka.partition.number: "3"
    kafka.replication.factor: "3"
spec:
  channeltypes:
    primes_ct1:
      meta:
        name: "primes_ct1"
      schema: '{"type":"int"}'

  channels:
    primes_ch1:
      meta:
        name: "primes_ch1"
        reference: ""
        Annotations: 
          kafka.partition.number: "3"
          kafka.replication.factor: "3"
        parent: ""  
      spec:
        type: "primes_ct1"
         
  apps:
    # number generators
    generator:
      meta:
        name: "generator"
      spec:
        node:
          meta:
            name: "node-generator"
            parent: "generator"
          spec:
            image: gcr.io/red-inspr/inspr/examples/primes/generator:latest
            replicas: 8
            environment:
              MODULE: 100
        boundary:
          input:
            - primes_ch1
          output:
            - primes_ch1

    # prints the filtered
    filter: 
      meta:
        name: "printer-primes"
        reference: ""
        parent: ""
      spec:        
        node:
          meta:
            name: "node-printer"
            parent: "printer"
          spec:
            image: gcr.io/red-inspr/inspr/examples/primes/printer:latest
            replicas: 2            
        boundary:
          input:
            - primes_ch1
          output:
            - primes_ch1
```
