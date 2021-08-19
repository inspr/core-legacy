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

```
$ helm install <your.release.name> inspr/insprd --set insprd.apps.createNamespace=true
```

Check your releases installed

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

