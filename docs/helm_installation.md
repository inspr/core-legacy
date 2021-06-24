# Helm Chart Installation

This document will walk you through the installation of Inspr using Helm, as well as showing the available configurations. Make sure you have [helm](https://helm.sh/) and [kubectl](https://kubernetes.io/docs/tasks/tools/) installed.

## Adding the Inspr Helm Repository

If you’re installing the Inspr chart via Helm, first you need to add the Inspr repository with the command:

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

The following table lists the configurable parameters of Inspr Helm Chart and their default values.

| Parameter | Description | Default
|--|--|--|
| replicas.insprdCount | Number of replicas of Insprd (Inspr daemon) deployment | 1 |
| replicas.authCount | Number of replicas of Authentication Service deployment | 1 |
| insprIngress.host | Main route for the Inspr Ingress Controller | inspr.com |
| image.insprdPullPolicy | Insprd's image pull policy | IfNotPresent |
| image.authPullPolicy | Auth Service image pull policy | IfNotPresent |
| insprEnvironment | Inspr Service Environment.  | test |
| deployment.insprdName | Name of Insprd deployment  | insprd-deployment |
| deployment.authName | Name of Auth Service deployment  | auth-deployment |
| deployment.initKey | Key used to initialize the Auth Service  | "1234567890" |
| insprdRepository | Insprd image  | gcr.io/red-inspr/insprd |
| authRepository | Auth Service image | gcr.io/red-inspr/authsvc |
| secretRepository | Secret Job Generator image  | gcr.io/red-inspr/secretgen |
| service.type | Kubernetes Service type | ClusterIP |
| service.authName | Auth Service Kubernetes Service name | insprd-svc |
| service.insprdName | Insprd Kubernetes Service name | insprd-svc |
| service.secretgenName | Secret Job Generator Kubernetes Service name | secretgen-svc |
| service.insprdPort | HTTP port of Insprd k8s service  | 80 |
| service.insprdTargetPort | Targeted port of insprdPort | 8080 |
| service.authPort | HTTP port of Auth Service k8s service  | 80 |
| service.authTargetPort | Targeted port of the authPort | 8080 |
| sidecarClient.readPort | Port which the Sidecar Client will receive requests | 3046 |
| sidecarClient.writePort | Port which the Load Balancer Sidecar will receive write requests from the Sidecar Client | 3048 |
| lbSidecar.port | Port which the Load Balancer Sidecar will receive read requests and redirect to the Sidecar Client | 3051 |
| lbSidecar.readPort | Port which the Load Balancer Sidecar is running | 3047 |
| lbsidecarImage | Load Balancer sidecar image | gcr.io/red-inspr/inspr/sidecar/lbsidecar |
| insprAppsNamespace | Kubernetes namespace on which Inspr apps will be instantiated | default |

## Exposing Inspr via NGINX

Insprd will be exposed via NGINX by default, so you just need to have NGINX installed. Follow the instructions in the [official NGINX documentation](https://kubernetes.github.io/ingress-nginx/deploy/) to install it.


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

