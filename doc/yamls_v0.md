
# Yamls Documentation

## DApps

| Field                     | Meaning   |
| ---                       | ---       |
| meta                      | metadata of DApp      |
| \|&rarr; name             | defines DApp name |
| \|&rarr; reference        | url to the inspr repository containing a already constructed DApp, it will load from this address the image containing all the necessary information for the creation of this DApp in your cluster.      |
| \|&rarr; parent           | defines DApp context in relation to the clust for example `*.app1.app2` would mean that this app is located on the path `root->app1->app2->app-name`. It is injected by the inspr daemon.    |
| \|&rarr; sha256           | tags images with their sha256 digest.     |
| spec                      | specification of DApp      |
| \|&rarr; Node             |       |
| \|&rarr; \|&rarr; Meta    |       |
| \|&rarr; \|&rarr; Spec    |       |
| \|&rarr; Apps             |       |    
| \|&rarr; Channels         |       |    
| \|&rarr; ChannelTypes     |       |    
| \|&rarr; Boundary         |       |
| \|&rarr; \|&rarr; Input   |       |
| \|&rarr; \|&rarr; Output  |       |    



## Channels

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
| \|&rarr; item_app_name        | name of the DApp currently using this channel     |


## Channel_Types

| Field                         | Meaning |
| ---                           | ---     |
| meta                          | metadata of Channel_Type     |
| \|&rarr; name                 | channel_type_name   |
| \|&rarr; reference            | url reference to the channel_type definition in the inspr repository, there are already well defined channel_types that can be used instead of defining your own.     |
| \|&rarr; parent               | It is injected by the inspr daemon and it's string composed of it's location's path, for example `'*.app1.app2'` means that the channel type belongs to the app2 in your cluster.       |
| \|&rarr; sha256               | tags images with their sha256 digest.     |
| schema                        | defines the message structure  that goes through this channel_type, example:  `'{"type":"int"}'`     |
| connectedchannels             | Is a list of channels that are created using this specific type, this is injected through the `inspr_cli`/ `inspr_daemon` |
| \|&rarr; item_channel         | name of the channel currently using this type     |



## General file

