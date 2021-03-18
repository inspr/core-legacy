# Yamls Documentation


Something in **markdown**.

<p>Then an HTML tag with crazy **markup** _all over_ the place!</p>



## DApps

| Field | Meaning   |
| --- | --- |
| meta                                      | metadata of DApp      |
| meta&rarr;name                             | defines DApp name |
| meta&rarr;reference                        | url to the inspr repository containing a already constructed DApp, it will load from this address the image containing all the necessary information for the creation of this DApp in your cluster.      |
| meta&rarr;parent                           | defines DApp context in relation to the clust for example `*.app1.app2` would mean that this app is located on the path `root->app1->app2->app-name`. It is injected by the inspr daemon.    |
| meta&rarr;sha256                           | tags images with their sha256 digest.     |
| spec                                      | specification of DApp      |
| meta&rarr;Node                             | Section describing the DApp node       |
| meta&rarr;Node&rarr;Meta                    |       |
| meta&rarr;Node&rarr;Meta&rarr; name           | defines node name |
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
{: .custom-class}


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
| \|&rarr; item_DApp            | name of the DApp currently using this channel     |


## Channel_Types

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


## Channel Type Structure


Channel_Type:
 - meta: <br>
    metadata of Channel_Type
    - name: <br>
        channel_type_name
    - reference: <br>
        url reference to the channel_type definition in the inspr repository, there are already well defined channel_types that can be used instead of defining your own. 
    - parent:<br>
        It is injected by the inspr daemon and it's string composed of it's location's path, for example `'*.app1.app2'` means that the channel type belongs to the app2 in your cluster.
    - sha256:<br>
        tags images with their sha256 digest.
 - schema: <br>
    defines the message structure  that goes through this channel_type, example:  `'{"type":"int"}'`
 - connected: <br>
    Is a list of channels that are created using this specific type, this is injected through the `inspr_cli`/ `inspr_daemon`
    - item_channel: <br>
    name of the channel currently using this type



## General file

