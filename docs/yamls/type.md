## Types 

### Definitions

| Field             | Meaning                                                                                                                                                            |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| apiVersion        | Specify what version of the API to use, for example `v1`                                                                                                         |
| kind              | Specifies which structure the file represents, in this case it would be `type`                                                                                     |
| meta              | Metadata of Type                                                                                                                                                   |
| &rarr;name        | Type's Name                                                                                                                                                        |
| &rarr;reference   | String that contains the url to the location of the Type definition in inspr's registry                                                                            |
| &rarr;annotations | Definitions that can describe characteristics of the Type that later on can be used to process/group the Types in your cluster.                                    |
| &rarr;parent      | Defines the Type context in the cluster through the path of the dApp in which it is stored, for example: `app1.app2` means that the Type is defined in the `app2`. |
|                   |
| schema            | defines the data structure that goes through this Type, example:  `'{"type":"int"}'`                                                                               |
| connectedchannels | Is a list of Channels names that are created using this specific type.                                                                                             |


### YAML example
```yaml
apiVersion: v1
kind: type
meta:
  name: primes_ct1  
schema: '{"type":"int"}'
```

[back](index.md)
