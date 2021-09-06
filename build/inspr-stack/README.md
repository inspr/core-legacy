# Inspr Stack helm chart

This chart creates a Inspr Stack deployment on a Kubernetes cluster using the Helm package manager, which consist of installing both Inspr Daemon and UIDP at the same time.


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
$ helm search repo inspr-stack

NAME                    CHART VERSION   APP VERSION     DESCRIPTION                                       
inspr/inspr-stack       v0.1.3          v0.1.3          Helm chart for installing Inspr's Deamon applic...
```

## Install Chart
---

To install the inspr-stack chart on helm run the following command(notice that the apps namespace is being setting to true, overwriting the chart value and enabling the app.namespace to be created):

```
$ helm install <your.release.name> inspr/inspr-stack --set insprd.apps.createNamespace=true
```

To check your releases installed:

```
$ helm list

NAME                NAMESPACE   REVISION    UPDATED                                 STATUS      CHART               APP VERSION
your.release.name   default     1           2021-08-19 10:55:57.498118205 -0300 -03 deployed    nspr-stack-v0.1.3   v0.1.3
```

## Uninstall Chart
---

```
$ helm uninstall <your.release.name>
```

<br>

# Usefull Overwrites  

## Example 1

How to enable the uidp ingress configuration

```
$ helm install <your.release.name> inspr/inspr-stack --set uidp.ingress.enabled=true
```

## Example 2

You can change the service port if the default one is already been used

```
$ helm install <your.release.name> inspr/inspr-stack \
--set insprd.service.port=<new.port> \
--set uidp.service.port=<new.port>

...

--set global.inspr.port=<new.port>

```

## Example 3

There are more log levels that you can use on the cluster besides the default one("info")

```
$ helm install <your.release.name> inspr/inspr-stack --set global.logLevel=<your.logLvl>
```

## Example 4

Changing the uidp admin password to your own password

```
$ helm install <your.release.name> inspr/inspr-stack --set uidp.admin.password=<your.password>
```

## Example yaml file overwrite

Using a yaml file to overwrite the char values

```
$ helm install <your.release.name> inspr/inspr-stack -f your_values.yaml
```

The example of a .yaml file that overwrite the values for the installation of the inspr-stack, which is contemplates both insprd and uidp subchart values.

```yaml
global:
  imagePullSecrets: []
  logLevel: info
  imageRegistry:
  insprd:
    name: insprd
    service:
      port: 80


insprd:
  enabled: true
  name: insprd
  image:
    registry: gcr.io/insprlabs
    repository: insprd
    tag: v0.1.3
  imagePullPolicy: IfNotPresent

  logLevel: info

  replicaCount: 1

  apps:
    createNamespace: true

  ingress:
    enabled: false
    host:
    class:

  initKey: ""

  service:
    type: ClusterIP
    port: 80
    targetPort: 8080

  sidecar:
    image: 
      registry: gcr.io/insprlabs
      repository: inspr/sidecar/lbsidecar
      tag: v0.1.3
    ports:
      client:
        read: 3046
        write: 3048
      server:
        read: 3047
        write: 3051

  auth:
    name: "auth"
    logLevel: info
    service:
      type: ClusterIP
      port: 80
      targetPort: 8081
    image:
      registry: gcr.io/insprlabs
      repository: authsvc
      tag: v0.1.3 

uidp:
  name: uidp
  enabled: true
  logLevel: info
  image:
    registry: gcr.io/insprlabs
    repository: uidp/redis/api
    tag: v0.1.3

  imagePullPolicy: IfNotPresent

  service:
    type: ClusterIP
    port: 80
    targetPort: 9001

  secret:
    name: '{{ .Release.Name }}-init-secret'
    image:
      registry: gcr.io/insprlabs
      repository: uidp/redis/secret
      tag: v0.1.3

  admin:
    password:
    token:
    generatePassword: true

  redis:
    create: true

  ingress:
    enabled: false
    class:
    host:

  insprd:
    init:
      enabled: true
      secret:
        key: key
        name: '{{ include "insprd.fullname" $ }}-init-key'
```

To see all the possible values overrides go to [Values_configuration](../../docs/values_configuration.md)