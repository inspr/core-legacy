
# Channels 

## Definitions

| Field             | Meaning                                                                                                                                                                  |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| apiVersion        | Specify what version of the API to use, for example `"v1"`                                                                                                               |
| kind              | Specifies which structure the file represents, in this case it would be `channel`                                                                                        |
| meta              | Metadata of Channel                                                                                                                                                      |
| &rarr; name       | Defines the Channel name                                                                                                                                                 |
| &rarr; reference  | String that is utilized to define certain tags to the Channel in question, a way for the user to categorize the numerous Channels in the cluster.                        |
| &rarr;Annotations | Definitions that can describe characteristics of the Channel that later on can be used to process/group the Channels in your cluster.                                    |
| &rarr; parent     | Defines the Channel context in the cluster through the path of the dApp in which it is stored, for example: `app1.app2` means that the Channel is defined in the `app2`. |
| &rarr; sha256     | Tags images with their sha256 digest.                                                                                                                                    |
| spec              |                                                                                                                                                                          |
| &rarr; type       | Defines the type of the Channel, this is a string that contains the name of any Type in the same context as the dApp that the channel is being created on.       |
| connectedapps     | List of dApp names that are using this Channel, this is injected by the Inspr daemon                                                                                     |
| connectedaliasses | A simple list of the aliases that are being used to reference this channel.                                                                                              |

## YAML example
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

[back](index.md)