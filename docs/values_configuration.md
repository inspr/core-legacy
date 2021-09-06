# Helm chart values Configuration

The following table lists the configurable parameters of Inspr Helm Chart and their default values. They can be setted with the flag on the helm command to overwrite the default ones.

## Insprd values

| Charts | Parameter | Description | Default
|--|--|--|--|
| insprd | global.imagePullPolicy | Insprd's image pull policy | IfNotPresent |
| insprd | name| Name of Insprd deployment | insprd |
| insprd | image.registry | Insprd image | gcr.io/insprlabs |
| insprd | image.repository | The name of the Docker image for the Insprd containers running | insprd |
| insprd | replicaCount | Number of replicas of Insprd (Inspr daemon) deployment | 1 |
| insprd | logLevel | Configures the log verbosity of the injector for Insprd | info & debug |
| insprd | apps.namespace | Kubernetes namespace on which Inspr apps will be instantiated | "{{ .Release.Name }}-inspr-apps" |
| insprd | apps.createNamespace | Allow to create a separated name space for the insprd apps | false |
| insprd | ingress.enable | If set to true, an Ingress service will be created. | false |
| insprd | ingress.host | Main route for the Inspr Ingress Controller | inspr-stack |  |
| insprd | ingress.class | Set the type that ingress will use |  |
| insprd | generateKey | If set to true, generate a insprd init key | true |
| insprd | key | Set the initKey value | "" |
| insprd | service.type | Sets the type of service to create for Insprd | ClusterIP |
| insprd | service.port | HTTP port of Insprd k8s service | 80 |
| insprd | service.targetPort | Targeted port of Insprd port | 8080 |
| insprd | sidecar.image.registry | Insprd's sidecar image | gcr.io/insprlabs |
| insprd | sidecar.image.repository | The name of the Docker image for the Sidecar containers running | inspr/sidecar/lbsidecar |
| insprd | sidecar.ports.client.read | Port which the Sidecar Client will receive requests | 3046 |
| insprd | sidecar.ports.client.write | Port which the Load Balancer Sidecar will receive write requests from the Sidecar Client | 3048 |
| insprd | sidecar.ports.server.read | Port which the Sidecar Server will receive requests | 3047 |
| insprd | sidecar.ports.server.write | Port which the Load Balancer Sidecar will receive write requests from the Sidecar Server | 3051 |
| insprd | auth.name | Name of Auth Service | auth |
| insprd | auth.service.type | Sets the type of service to create for Auth Service | ClusterIP |
| insprd | auth.service.port | HTTP port of Auth Service k8s service | 80 |
| insprd | auth.service.targetPort | Targeted port of the auth port | 8081 |
| insprd | auth.image.registry | Auth Service image | gcr.io/insprlabs |
| insprd | auth.image.repository | The name of the Docker image for the Auth Service containers running | authsvc |
| insprd | secretGenerator.image.registry | Secret Generator image | gcr.io/insprlabs |
| insprd | secretGenerator.image.repository | The name of the Docker image for the Secret Generator containers running | secretgen |

## UIDP values

| Charts | Parameter | Description | Default
|-|-|-|-|
| uidp | global.redis.password | Set the global redis password | |
| uidp | name| Name of UIDP deployment | uidp |
| uidp | image.registry | UIDP image | gcr.io/insprlabs |
| uidp | image.repository | The name of the Docker image for the UIDP containers running | uidp/redis/api |
| uidp | imagePullPolicy | UIDP's image pull policy | IfNotPresent |
| uidp | logLevel | Configures the log verbosity of the injector for UIDP | info & debug |
| uidp | service.type | Sets the type of service to create for UIDP | ClusterIP |
| uidp | service.port | HTTP port of UIDP k8s service | 80 |
| uidp | service.targetPort | Targeted port of UIDP port | 9001 |
| uidp | secret.name | Name of UIDP Secret | {{ .Release.Name }}-init-secret |
| uidp | secret.image.registry | UIDP secret image | gcr.io/insprlabs |
| uidp | secret.image.repository | The name of the Docker image for the UIDP secret containers running | uidp/redis/api |
| uidp | admin.generatePassword | If set to true, generate a admin password | true |
| uidp | admin.password | Set the admin password |  |
| uidp | redis.create | Initialise | true |
| uidp | ingress.enabled | If set to true, an Ingress service will be created. | false |
| uidp | ingress.class | Set the type that ingress will use |  |
| uidp | ingress.host | Sets the hostname for the Route |  |
| uidp | insprd.name | Name of Insprd deployment | insprd |
| uidp | insprd.init.enabled | If set to true, an Insprd service will be created | false |
| uidp | insprd.init.secret.key | The secret key value |  |
| uidp | insprd.init.secret.name | Name of the key where the requested secret value is located in the Kubernetes secret |  |
| uidp | insprd.token | Set the Insprd token |  |

## Inspr Stack values

When trying to use the inspr-stack chart we can define their values in a similar way to the insprd and UIDP, the difference is the necessity of adding their dependency before the parameter for the value that you want to change.

Insprd example:
```yaml
...
insprd: # <= dependency
  name: insprd # <= parameter
...
```  

UIDP example:
```yaml
...
uidp: # <= dependency
  name: uidp # <= parameter
...
```

| Charts | Parameter | Description | Default
|--|--|--|--|
| inspr-stack | insprd.name | Name of insprd deployment | insprd |
| inspr-stack | uidp.name | Name of UIDP deployment | uidp |