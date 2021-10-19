# inspr UIDP helm chart

This chart creates a Inspr UIDP deployment on a Kubernetes cluster using the Helm package manager.

## Prerequisites
---
- Kubernetes 1.16+
- Helm 3+

## Get Repo Info
---
If you’re installing the Inspr chart via Helm, first you need to add the Inspr repository with the command:

```
$ helm repo add inspr https://inspr-charts.storage.googleapis.com/
```
Then, make sure everything is up-to-date:

```
$ helm repo update
```

And finaly to see the details from Inspr on the helm you can do:

```
$ helm search repo uidp

NAME                    CHART VERSION   APP VERSION     DESCRIPTION                                       
inspr/uidp              v0.1.3          v0.1.3          A Helm chart for Kubernetes 
```

## Install Chart
---

To install the uidp chart on helm run the following command:

```
$ helm install <your.release.name> inspr/uidp
```

To check your releases installed:

```
$ helm list

NAME                NAMESPACE   REVISION    UPDATED                                 STATUS      CHART           APP VERSION
your.release.name   default     1           2021-08-19 10:55:57.498118205 -0300 -03 deployed    uidp-v0.1.3     v0.1.3    
```

## Uninstall Chart
---

```
$ helm uninstall <your.release.name>
```

# Usefull Overwrites  

Threre are more than one way to overwrite the chart values, you can set value on the command line or point to a overwrite_value file.

Exemple command line with --set

```
$ helm install <your.release.name> inspr/inspr-stack --set <parameter>=<value>
```

Exemple command line pointing to a file

```
$ helm install <your.release.name> inspr/inspr-stack -f your_values.yaml
```
# Usefull Overwrites  

Threre are more than one way to overwrite the chart values, you can set value on the command line or point to a overwrite_value file.

Exemple command line with --set

```
$ helm install <your.release.name> inspr/inspr-stack --set <parameter>=<value>
```

## Example yaml file overwrite

Exemple command line pointing to a file

```
$ helm install <your.release.name> inspr/inspr-stack -f your_values.yaml
```

Example of a .yaml file that overwrite the values for the installation of the uidp chart.

```yaml
global:
  redis:
    password:

name: uidp

image:
  registry: gcr.io/insprlabs
  repository: uidp/redis/api
  tag: v0.1.4

imagePullPolicy: IfNotPresent

logLevel: info

service:
  type: ClusterIP
  port: 80
  targetPort: 9001

secret:
  name: '{{ .Release.Name }}-init-secret'
  image:
    registry: gcr.io/insprlabs
    repository: uidp/redis/secret
    tag: v0.1.4

admin:
  generatePassword: true
  password:

redis:
  create: true
  fullNameOverride: redis
  password:
  existingSecret: "uidp-redis-secret"
  existing:
    host:
    port:
  cluster:
    nodes: 3

ingress:
  enabled: false
  class: "nginx"
  host:

insprd:
  name: "insprd"
  init:
    enabled: false
    secret:
      key:
      name:
  address: 'insprd'
  token:
```

To see usefull overwrites go to [Values_configuration](../../docs/values_configuration.md)