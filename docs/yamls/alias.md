## Alias 

### Definitions

| Field             | Meaning                                                                                                                                                                                                      |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| apiVersion        | Specify what version of the API to use, for example `"v1"`                                                                                                                                                   |
| kind              | Specifies which structure the file represents, in this case it would be `alias`                                                                                                                              |
| meta              | Metadata of Type                                                                                                                                                                                             |
| &rarr;name        | Type Name of the Alias                                                                                                                                                                                       |
| &rarr;reference   | String that is utilized to define certain tags to the Alias in question, a way for the user to categorize the numerous Alias in the cluster.                                                                 |
| &rarr;Annotations | Definitions that can describe characteristics of the Alias that later on can be used to process/group the Alias in your cluster.                                                                             |
| &rarr;parent      | Defines the Alias context in the cluster through the path of the dApp in which it is stored, for example: `app1.app2` means that the Alias is defined in the `app2` and that this DApp is a child of `app1`. |
|                   |
| &rarr;UUID        | String that is a unique ID,, it is filled automatically when creating certain components. By default leave it empty                                                                                          |
| Target            | specifies the name of the Channel in which this Alias points to.                                                                                                                                             |


### YAML example
```yaml
kind: alias
apiVersion: v1
meta:
  name: aliasOne
  parent: chcheck
target: testch0
```

[back](index.md)
