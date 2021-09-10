
# Steps to install the external dns in the cluster

[External-DNS](https://github.com/kubernetes-sigs/external-dns) is a tool that
allows us to set a deployment in kubernetes that watches services and ingresses
in your kubernetes cluster with the purporse of associating a subdomain to them.

### Requirements

For the usage of the commands in this documentation you will need:
- helm
- domain registered in a supported domain registry, list [here](https://github.com/kubernetes-sigs/external-dns#roadmap)


### Installation

For this to work in a proper manner we need to tell the external-dns deployment
what is the domain that we own and what is our project in which the cluster is
located. In this documentation we are using gcloud dns provider and we own the
inspr.dev domain, therefore the examples will be either using these values or
`<clusterDomain>` to reference the url owned.

Before actually specifying the values for the external-dns deployment we need to
interact with the cloud provider and create a service.account that allows the
application to interact with the DNS records of the project. Taking as a base
the tutorial found
[here](http://tech.paulcz.net/kubernetes-cookbook/gcp/gcp-external-dns/) we
decided to setup the permissions as a k8s secrets that is stored in the cluster.

The set of instructions that create the secret containing the values are:
```
$ gcloud iam service-accounts create external-dns \
    --display-name "Service account for ExternalDNS on GCP"

$ gcloud projects add-iam-policy-binding <GCP_PROJECT_ID> \
    --role='roles/dns.admin' \
    --member='serviceAccount:external-dns@<GCP_PROJECT_ID>.iam.gserviceaccount.com'
...

$ gcloud iam service-accounts keys create credentials.json \
    --iam-account external-dns@<GCP_PROJECT_ID>.iam.gserviceaccount.com
```

From the `credential.json` file create we can deploy a k8s secret that stores
the information in a safe manner.

```
$ kubectl -n external-dns-gcp create secret \
    generic external-dns \
  --from-file=credentials.json=credentials.json
```

**Remind yourself to delete the `credential.json` and never store it in a public
place, it contains the permissions to interact with the dns records of your
project.**

A pratical way to install the external-dns is to use a yaml file to specify the
changes that you want to enforce in the helm chart installation. 

There are two links that provides a description of the exitent values of the
external-dns helm chart:
- [github](https://github.com/bitnami/charts/blob/master/bitnami/external-dns/values.yaml)
- [artifachub](https://artifacthub.io/packages/helm/bitnami/external-dns)
Both of them can be used as a base by anyone, in our case we want only a small amount
of changes so the yaml content below should suffice for the installation in the
development environment.


```yaml
# values.yaml
# contains external-dns config

# allows to delete the dns registry when deleting a service
policy: upsert-only

# specifies the suffix, only will act on ingresses and services that have their
# host/annotation ending with this domain
clusterDomain: inspr.dev

# cloud provider configuration
provider: google
google: 
    project: insprlabs
    serviceAccountSecret: external-dns
```

With the file defined we can just run the following command
```bash
$ helm install my-release -f values.yaml bitnami/external-dns
```


After the installation the external-dns should be up and running in your
cluster, therefore we only need to deploy a service with a annotation specifying
the url or a ingress with a host ending with the value of `clusterDomain`

#### kubectl

If you prefer to use kubectl, all the necessary instructions and steps can be
found
[here](https://github.com/kubernetes-sigs/external-dns/tree/master/docs/tutorials).

## Usage Examples

#### Simple service with annotation
```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    external-dns.alpha.kubernetes.io/hostname: nginx.external-dns-test.<cluster_domain>
  name: nginx
spec:
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: nginx
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
        ports:
        - containerPort: 80
```
#### Ingress with hostname
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx
spec:
  rules:
  - host: via-ingress.external-dns-test.<cluster_domain>
    http:
      paths:
      - backend:
          serviceName: nginx
          servicePort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: nginx
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
        ports:
        - containerPort: 80
```


