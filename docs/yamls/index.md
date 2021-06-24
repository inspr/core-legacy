# YAML Documentation

The creation of YAML files allows the proper usage of the Inspr CLI.

In this way, when the file is written in one of the formats described below it can be processed and applied to the cluster through the CLI commands `insprctl apply -f <file_path>` or `insprctl apply -k <files_dir>`.


## dApps
> A dApp is an Inspr structure that allows the user to contain his application in a component that can interact with the other dApps in the cluster. In essence is a layer of abstration that facilitates the communication between multiple microservices by layerying each of one of them in a dApp.

[definitions and examples](dapp.md)

## Channels
> Responsible for creating a message broker's topic, a Channel serve as a path in which two or more dApps can exchange data. It must have a Type defined.

[definitions and examples](channel.md)

## Types
> Responsible for defining the message format for any Channel defined with this Type.
> 
> A Type must always have it's schema specified, this has to be an [avro structure](https://avro.apache.org/docs/current/).
> This would be either a string containing a json structure specifying the format of the message or a path to a file containing such information.

[definitions and examples](type.md)

## General file

>The so called general file, or composed file, is nothing more than a YAML that congregates two or more definitions of the elements described above into a single dApp. 

For example a basic example dApp, that has a collection of other smaller dApps like number-generator and filter, plus some definitions of Types and Channels.

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
  Types:
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
