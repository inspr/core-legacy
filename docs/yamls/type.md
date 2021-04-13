## Channel Types 

### Definitions

| Field             | Meaning                                                                                                                                                                            |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| apiVersion        | Specify what version of the API to use, for example `"v1"`                                                                                                                         |
| kind              | Specifies which structure the file represents, in this case it would be `channeltype`                                                                                              |
| meta              | Metadata of Channel Type                                                                                                                                                           |
| &rarr;name        | Channel Type Name                                                                                                                                                                  |
| &rarr;reference   | String that is utilized to define certain tags to the Channel Type in question, a way for the user to categorize the numerous Channel Types in the cluster.                        |
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

[back](index.md)