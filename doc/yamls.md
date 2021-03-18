# Yamls Documentation

The creation of yaml files allows for the proper usage of the Inspr cli.

The reason for this is that the file when written in a one of the formats described bellow it can be processed and applied to the cluster throught the `inspr apply -f <file_path>` or `inspr apply -k <files_dir>` commands.

## DApps

### Definitions

| Field | Meaning   |
| --- | --- |
| meta                                      | metadata of DApp      |
| \|&rarr;name                             | defines DApp name |
| \|&rarr;reference                        | url to the inspr repository containing a already constructed DApp, it will load from this address the image containing all the necessary information for the creation of this DApp in your cluster.      |
| \|&rarr;parent                           | defines DApp context in relation to the clust for example `*.app1.app2` would mean that this app is located on the path `root->app1->app2->app-name`. It is injected by the inspr daemon.    |
| \|&rarr;sha256                           | tags images with their sha256 digest.     |
| spec                                      | specification of DApp      |
| \|&rarr;Node                             | Section describing the DApp node       |
| \|&rarr;\|&rarr;Meta                    |       |
| \|&rarr;\|&rarr;\|&rarr; name           | defines node name |
| \|&rarr; \|&rarr; \|&rarr; reference      | url to the inspr repository containing a already constructed DApp-Node, it will load from this address the image containing all the necessary information for the creation of this node in your cluster.      |
| \|&rarr; \|&rarr; \|&rarr; parent         | defines the node context in relation to the clust for example `*.app1.app2` would mean that this node is located on the path `root->app1->app2`. It is injected by the inspr daemon.    |
| \|&rarr; \|&rarr; \|&rarr; sha256         | tags images with their sha256 digest.     |
| \|&rarr; \|&rarr; Spec                    |       |
| \|&rarr; \|&rarr; \|&rarr; Image          | url to the location of the already defined node in the inspr repository      |
| \|&rarr; \|&rarr; \|&rarr; Replicas       | defines the amount of replicas to be created in your cluster       |
| \|&rarr; \|&rarr; \|&rarr; Envioronment   | defines the envioronment variables of your pods      |
| \|&rarr; Apps                             | set of DApps that are connected to this DApp, can be either specified when creating a new app or is modified by the inspr daemon when creating/updating different DApps      |    
| \|&rarr; Channels                         | set of Channels that are created in the context of this DApp      |    
| \|&rarr; ChannelTypes                     | set of Channel Types that are created in the context of this DApp      |    
| \|&rarr; Boundary                         |       |
| \|&rarr; \|&rarr; Input                   | List of channels that are used for the input of this DApp      |
| \|&rarr; \|&rarr; Output                  | List of channels that are used for the output of this DApp      |    

### yaml example
```
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

| Field                         | Meaning |
| ---                           | ---     |
| meta                          | metadata of Channel     |
| \|&rarr; name                 | defines the Channel name     |
| \|&rarr; reference            | url reference to the channel definition in the inspr repository, there are already well defined channel that can be used instead of defining your own.     |
| \|&rarr; parent               | it is injected by the inspr daemon, defines the Channel context in the cluster through the path of the app in which the channel is stored, for example: "*.app1.app2" means that the channel is defined in the app2.    |
| \|&rarr; sha256               | tags images with their sha256 digest.     |
| spec                          |      |
| \|&rarr; type                 | defines the type of the channel, this field is a string that contains the name of any of the channel_types defined in your cluster     |
| connectedapss                 |  List of app names that are using this channel, this is injected by the inspr daemon    |
| \|&rarr; item_DApp            | name of the DApp currently using this channel     |

### yaml example
```
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


## Channel_Types 

### Definitions

| Field                         | Meaning |
| ---                           | ---     |
| meta                          | metadata of Channel_Type     |
| \|&rarr;name                 | channel_type_name   |
| \|&rarr;reference            | url reference to the channel_type definition in the inspr repository, there are already well defined channel_types that can be used instead of defining your own.     |
| \|&rarr;parent               | It is injected by the inspr daemon and it's string composed of it's location's path, for example `'*.app1.app2'` means that the channel type belongs to the app2 in your cluster.       |
| \|&rarr;sha256               | tags images with their sha256 digest.     |
| schema                        | defines the message structure  that goes through this channel_type, example:  `'{"type":"int"}'`     |
| connectedchannels             | Is a list of channels that are created using this specific type, this is injected through the `inspr_cli`/ `inspr_daemon` |
| \|&rarr;item_channel         | name of the channel currently using this type     |


### list definitions
Channel_Type:
 - meta: 
    - name: channel_type_name
    - reference: url reference to the channel_type definition in the inspr repository, there are already well defined channel_types that can be used instead of defining your own. 
    - parent:It is injected by the inspr daemon and it's string composed of it's location's path, for example `'*.app1.app2'` means that the channel type belongs to the app2 in your cluster.
    - sha256: tags images with their sha256 digest.
 - schema: defines the message structure  that goes through this channel_type, example:  `'{"type":"int"}'`
 - connected: Is a list of channels that are created using this specific type, this is injected through the `inspr_cli`/ `inspr_daemon`
    - item_channel: name of the channel currently using this type

### yaml example
```
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

The so called general file or composed file is nothing more than a yaml that congregates two or more definitions of one of the elements above. 

For example an App that has a collection of other apps plus some definitions of channel_types and channels.

### yaml example

```
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
    primes_ct2:
      meta:
        name: "primes_ct2"
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
    primes_ch2:
      meta:
        name: "primes_ch2"
        reference: ""
        Annotations:     
          kafka.partition.number: "3"
          kafka.replication.factor: "3"
        parent: ""  
      spec:
        type: "primes_ct2"
        
         
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
            - primes_ch2
          output:
            - primes_ch1

    # filters primes
    filter: 
      meta:
        name: "filter-primes"
        reference: ""
        parent: ""
      spec:
        node:
          meta:
            name: "node-filter"
            parent: "filter-primes"
          spec:
            image: gcr.io/red-inspr/inspr/examples/primes/filter:latest
            replicas: 2            

        boundary:
          input:
            - primes_ch1
          output:
            - primes_ch2
        
        

    # prints the filtered
    printer: 
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
            - primes_ch2
          output:
            - primes_ch2
```
