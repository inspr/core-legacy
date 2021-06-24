
# Channels 

## Definitions

| Field              | Meaning                                                                                                                                                                                                                                    |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| apiVersion         | Specify what version of the API to use, for example `"v1"`                                                                                                                                                                                 |
| kind               | Specifies which structure the file represents, in this case it would be `channel`                                                                                                                                                          |
| meta               | Metadata of Channel                                                                                                                                                                                                                        |
| &rarr; name        | Defines the Channel name                                                                                                                                                                                                                   |
| &rarr; reference   | String that contains the url to the location of the Channel definition in inspr's registry                                                                                                                                                 |
| &rarr; annotations | Definitions that can describe characteristics of the Channel that later on can be used to process/group the Channels in your cluster.                                                                                                      |
| &rarr; parent      | Defines the Channel context in the cluster through the path of the dApp in which it is stored, for example: `app1.app2` means that the Channel is defined in the `app2`.                                                                   |
| spec               |                                                                                                                                                                                                                                            |
| &rarr; type        | This field is reponsible for the definition of the what type of message will be send through the channel, the content is a string the represents the name of a inspr structure called Type that has the `avro` definitions of the message. |
| connectedapps      | List of dApp names that are using this Channel, this is injected by the Inspr daemon                                                                                                                                                       |
| connectedaliasses  | A simple list of the aliases that are being used to reference this channel.                                                                                                                                                                |

## YAML example
```yaml
apiVersion: v1
kind: channel

meta:
  name: "channel_1"
  parent: ""
spec:
  type: "mbct1"
  brokerlist:
    - "kafka"

```

[back](index.md)