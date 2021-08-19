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

```
$ helm install <your.release.name> inspr/inspr-stack --set insprd.apps.createNamespace=true
```

Check your releases installed

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