## Alias 

### Definitions

| Field             | Meaning                                                                                                                                                                                                      |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| apiVersion        | Specify what version of the API to use, for example `v1`                                                                                                                                                   |
| kind              | Specifies which structure the file represents, in this case it would be `alias`                                                                                                                              |
| meta              | Metadata of Alias                                                                                                                                                                                            |
| &rarr;name        | Name of the Alias, it follows the format of `childDApp.childDAppBoundary`                                                                                                                                    |
| &rarr;annotations | Definitions that can describe characteristics of the Alias that later on can be used to process/group the Alias in your cluster.                                                                             |
| &rarr;parent      | Defines the Alias context in the cluster through the path of the dApp in which it is stored, for example: `app1.app2` means that the Alias is defined in the `app2` and that this DApp is a child of `app1`. |
|                   |
| target            | specifies the name of the Channel in which this Alias points to.                                                                                                                                             |


### YAML example
```yaml
kind: alias
apiVersion: v1
meta:
  name: "<dapp_name>.<boundary_name>"
  parent: <path_to_current_dapp>
target: <channel_name>
```

[back](index.md)
