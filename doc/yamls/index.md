# YAML Documentation

The creation of YAML files allows the proper usage of the Inspr CLI.

In this way, when the file is written in one of the formats described below it can be processed and applied to the cluster through the CLI commands `inspr apply -f <file_path>` or `inspr apply -k <files_dir>`.


## dApps
> A dApps is an Inspr structure that allows the user to contain his application in a container that can interact with the other components of the Inspr structure. One could call it an extra layer of abstraction that allows the simplification of communications between microservices in your cluster.

[definitions and examples](dapp.md)

## Channels
> Channels are a Inspr definition that facilitates the user's control over the message broker used to send message between dApps. 

[definitions and examples](channel.md)

## Channel Types
> Responsible for defining the message format for any channel defined with this type.
> 
> A Channel Type must always have it's schema specified, this would be a either a string containing a json structure specifing the format of the message or a path to a file containing such information.

[definitions and examples](type.md)

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
