# Helm Chart Installation

This document will walk you through the installation of Inspr using Helm, as well as showing the available configurations. Make sure you have [helm](https://helm.sh/) and [kubectl](https://kubernetes.io/docs/tasks/tools/) installed.

## Adding the Inspr Helm Repository

If you’re installing the Inspr chart via helm, first you need to add the Inspr repository with the command:

```
$ helm repo add inspr https://inspr-charts.storage.googleapis.com/
```

Then, make sure everything is up-to-date:

```
$ helm repo update
```
And finally, install it by running the command:

```
$ helm upgrade --install insprd inspr/insprd
```

## Configuration

The following table lists the configurable parameters of the Inspr Ingress controller chart and their default values.

| Parameter | Description | Default
|--|--|--|
| replicaCount | Number of replicas of the Inspr Ingress controller deployment. | 1 |
| insprIngress.host | Main route for the Inspr Ingress Controller | inspr.com |
| deployment.insprdName | Name of the Ingress Inspr deployment  | insprd-deployment|
| insprEnvironment | Inspr Service Environment.  | test |
| service.type | Insprd Kubernetes Service type | ClusterIP |
| service.insprdName | Insprd Kubernetes Service name | insprd-svc |
| service.insprdPort | HTTP port of the Inspr controller service.  | 80 |
| service.insprdTargetPort | Target port of the insprdPort. | 8080 |
| kafkaSidecarImage | Kafka operator sidecar image | gcr.io/red-inspr/inspr/sidecar/kafka |
| kafkaBootstrap | Kafka operator bootstrap configuration | kafka.default.svc:9092 |
| kafkaAutoOffsetReset | Kafka operator offset reset  | earliest |
| insprAppsNamespace | Kubernetes namespace on which Inspr apps will be instantiated | default |

## Exposing Inspr via NGINX

Inspr Controller will be exposed via NGINX by default, so you just need to have NGINX installed. Follow the instructions in the [official NGINX documentation](https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-helm/) to install it.


### Hostname

To properly access the cluster via the Inspr CLI it is necessary to configure the hostname in your machine to comport the Ingress host.

This can be done in unix machines editting the file `etc/hosts` and adding `<cluster_IP> <insprIngress.host>`. For further understanding of what is happening access this [link](https://debian-handbook.info/browse/stable/sect.hostname-name-service.html)

For Windows is recommended to follow the steps in [here](https://docs.microsoft.com/en-us/windows-server/networking/technologies/ipam/add-a-dns-resource-record).

## Exposing Inspr via port forward

After installing Inspr, you can also expose the port without using NGINX. To do so, follow the commands below:

First, get the `POD_NAME` with:
```
$ export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=insprd,app.kubernetes.io/instance=insprd" -o jsonpath="{.items[0].metadata.name}")
```
Then, the `CONTAINER_PORT`:
```
$ export CONTAINER_PORT=$(kubectl get pod --namespace default $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
```
And now, use both to expose Inspr via port forward:
```
$ kubectl --namespace default port-forward $POD_NAME 8080:$CONTAINER_PORT
```


## Installing Kafka via Helm

To install Kafka via Helm, first add the `bitnami` repository:

```
$ helm repo add bitnami https://charts.bitnami.com/bitnami
```

And then, install Kafka with:

```
$ helm install my-release bitnami/kafka
```

You can also check how to do the installation on the [official bitnami Kafka page](https://bitnami.com/stack/kafka/helm).


[Get Started](readme.md)

