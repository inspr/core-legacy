# Inspr Daemon Helm Chart

This chart creates a Inspr Deamon deployment on a Kubernetes cluster using the Helm package manager.

## Prerequisites
---

- Kubernetes 1.19+
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
$ helm search repo insprd

NAME                    CHART VERSION   APP VERSION     DESCRIPTION                                       
inspr/insprd            v0.1.3          v0.1.3          Helm chart for installing Inspr's Deamon applic...
```

## Install Chart
---

To install the insprd chart on helm run the following command(notice that the apps name space is being setting to true, overwriting the chart value and enabling the app.namespace to be created):

```
$ helm install <your.release.name> inspr/insprd --set apps.createNamespace=true
```

To check your releases installed:

```
$ helm list

NAME                NAMESPACE   REVISION    UPDATED                                 STATUS      CHART           APP VERSION
your.release.name   default     1           2021-08-19 10:50:33.742251146 -0300 -03 deployed    insprd-v0.1.3   v0.1.3
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

## Example yaml file overwrite

Exemple command line pointing to a file

```
$ helm install <your.release.name> inspr/inspr-stack -f your_values.yaml
```

Example of a .yaml file that overwrite the values for the installation of the insprd chart.

```yaml
global:
  imagePullSecrets: []
  imageRegistry:

name: "insprd"
image:
  registry: gcr.io/insprlabs
  repository: insprd
  tag: v0.1.3
imagePullPolicy: IfNotPresent

replicaCount: 1
logLevel: info

apps:
  namespace: "{{ .Release.Name }}-inspr-apps"
  createNamespace: false

ingress:
  enabled: false
  host:
  class:

init:
  generateKey: true
  key: ""

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
  service:
    type: ClusterIP
    port: 80
    targetPort: 8081
  image:
    registry: gcr.io/insprlabs
    repository: authsvc
    tag: v0.1.3 

secretGenerator:
  image:
    registry: gcr.io/insprlabs
    repository: secretgen
    tag: v0.1.3
```

To see usefull overwrites go to [Values_configuration](../../docs/values_configuration.md)