# YAML Documentation

The creation of YAML files allows the proper usage of the Inspr CLI.

In this way, when the file is written in one of the formats described below it can be processed and applied to the cluster through the CLI commands `inspr apply -f <file_path>` or `inspr apply -k <files_dir>`.

## dApps

### Definitions

| Field                | Meaning                                                                                                                                                                                     |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion           | Specify what version of the API to use, for example `v1`                                                                                                                                    |
| kind                 | Specifies which structure the file represents, in this case it would be `dapp`                                                                                                              |
| meta                 | Metadata of dApp                                                                                                                                                                            |
| &rarr;name           | Defines dApp name                                                                                                                                                                           |
| &rarr;reference      | String that is utilized to defined certain tags to the dApp in question, a way for the user to categorize the numerous dApps in the cluster.                                                |
| &rarr;Annotations    | Definitions that can describe characteristics of the dApp that later on can be used to process/group the dApps in your cluster.                                                             |
| &rarr;parent         | Defines dApp context in relation to the cluster. For example `app1.app2` would mean that this dApp is located on the path `root->app1->app2->app-name`. It is injected by the Inspr daemon. |
| &rarr;sha256         | Tags images with their sha256 digest.                                                                                                                                                       |
| spec                 | specification of dApp                                                                                                                                                                       |
| &rarr;Node           | Section describing the dApp Node                                                                                                                                                            |
| &rarr;&rarr;Meta     | Metadata of the Node                                                                                                                                                                        |
| name                 | Defines Node name                                                                                                                                                                           |
| reference            | String that is utilized to defined certain tags to the Node in question, a way for the user to categorize the numerous Nodes in the cluster.                                                |
| Annotations          | Definitions that can describe characteristics of the Node that later on can be used to process/group the Nodes in your cluster.                                                             |
| parent               | Defines the Node context in relation to the cluster for example `app1.app2` would mean that this Node is located on the path `root->app1->app2`. It is injected by the Inspr daemon.        |
| sha256               | tags images with their sha256 digest.                                                                                                                                                       |
| &rarr;&rarr; Spec    |                                                                                                                                                                                             |
| Image                | An URL that serve to point to the location in which the docker image of your application is stored                                                                                          |
| Replicas             | Defines the amount of replicas to be created in your cluster                                                                                                                                |
| Environment          | Defines the environment variables of your Node                                                                                                                                              |
| &rarr; Apps          | Set of dApps that are connected to this dApp, can be specified when creating a new dApp or modified when a dApp is updated.                                                                 |
| &rarr; Channels      | Set of Channels that are created in the context of this dApp                                                                                                                                |
| &rarr; ChannelTypes  | Set of Channel Types that are created in the context of this dApp                                                                                                                           |
| &rarr; Boundary      |                                                                                                                                                                                             |
| &rarr; &rarr; Input  | List of Channels that are used for the input of this dApp                                                                                                                                   |
| &rarr; &rarr; Output | List of Channels that are used as the output of this dApp                                                                                                                                   |

### YAML example
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

## Channels 

### Definitions

| Field             | Meaning                                                                                                                                                                  |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| apiVersion        | Specify what version of the API to use, for example `"v1"`                                                                                                               |
| kind              | Specifies which structure the file represents, in this case it would be `channel`                                                                                        |
| meta              | Metadata of Channel                                                                                                                                                      |
| &rarr; name       | Defines the Channel name                                                                                                                                                 |
| &rarr; reference  | String that is utilized to defined certain tags to the Channel in question, a way for the user to categorize the numerous Channels in the cluster.                       |
| &rarr;Annotations | Definitions that can describe characteristics of the Channel that later on can be used to process/group the Channels in your cluster.                                    |
| &rarr; parent     | Defines the Channel context in the cluster through the path of the dApp in which it is stored, for example: `app1.app2` means that the Channel is defined in the `app2`. |
| &rarr; sha256     | Tags images with their sha256 digest.                                                                                                                                    |
| spec              |                                                                                                                                                                          |
| &rarr; type       | Defines the type of the Channel, this is a string that contains the name of any Channel Type on the same context as the dApp that the channel is being created on.       |
| connectedapps     | List of dApp names that are using this Channel, this is injected by the Inspr daemon                                                                                     |
| connectedaliasses | A simple list of the aliases that are being used for to reference this channel.                                                                                          |

### YAML example
```yaml
apiVersion: v1
kind: channel
meta:
  name: primes_ch1  
  Annotations: 
    kafka.partition.number: 1
    kafka.replication.factor: 1  
spec:
  type: primes_ct1
```


## Channel Types 

### Definitions

| Field             | Meaning                                                                                                                                                                            |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion        | Specify what version of the API to use, for example `"v1"`                                                                                                                         |
| kind              | Specifies which structure the file represents, in this case it would be `channeltype`                                                                                              |
| meta              | Metadata of Channel Type                                                                                                                                                           |
| &rarr;name        | Channel Type Name                                                                                                                                                                  |
| &rarr;reference   | String that is utilized to defined certain tags to the Channel Type in question, a way for the user to categorize the numerous Channel Types in the cluster.                       |
| &rarr;Annotations | Definitions that can describe characteristics of the Channel Type that later on can be used to process/group the Channel Types in your cluster.                                    |
| &rarr;parent      | Defines the Channel Type context in the cluster through the path of the dApp in which it is stored, for example: `app1.app2` means that the Channel Type is defined in the `app2`. |
|                   |
| &rarr;sha256      | Tags images with their sha256 digest.                                                                                                                                              |
| schema            | defines the data structure that goes through this Channel Type, example:  `'{"type":"int"}'`                                                                                       |
| connectedchannels | Is a list of Channels names that are created using this specific type.                                                                                                             |


### YAML example
```yaml
apiVersion: v1
kind: channeltype
meta:
  name: primes_ct1  
schema: '{"type":"int"}'
```

## General file

### Definition

The so called general file, or composed file, is nothing more than a YAML that congregates two or more definitions of the elements described above into a single dApp. 

For example a basic example dApp, that has a collection of other smaller dApps like number-generator and filter, plus some definitions of channel Types and channels.

### YAML example

```yaml
apiVersion: v1
kind: dapp
meta:
  name: basic-example  
  Annotations: 
    kafka.partition.number: 3
    kafka.replication.factor: 3
spec:
  channeltypes:
    primes_ct1:
      meta:
        name: primes_ct1
      schema: '{"type":"int"}'

  channels:
    primes_ch1:
      meta:
        name: primes_ch1        
        Annotations: 
          kafka.partition.number: 3
          kafka.replication.factor: 3        
      spec:
        type: primes_ct1
         
  apps:
    # number generators
    generator:
      meta:
        name: generator
      spec:
        node:
          meta:
            name: node-generator
            parent: generator
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
        name: printer-primes
      spec:        
        node:
          meta:
            name: node-printer
            parent: printer
          spec:
            image: gcr.io/red-inspr/inspr/examples/primes/printer:latest
            replicas: 2            
        boundary:
          input:
            - primes_ch1
          output:
            - primes_ch1
```
