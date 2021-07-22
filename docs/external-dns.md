
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


When using helm to install, these settings can be passed via the `--set` flag
like the example below:

```bash
$ helm repo add bitnami https://charts.bitnami.com/bitnami

$ helm install my-release \
--set provider=google \
--set clusterdomain=inspr.dev \
--set policy=sync \
--set google.project=insprlabs \
bitnami/external-dns
```

Another way which is a bit more pratical is to use a yaml file to specify the
changes that you want to enforce in the external-dns installation. The
repository for the helm chart offers a [yaml
file](https://github.com/bitnami/charts/blob/master/bitnami/external-dns/values.yaml)
which can be used as a base by anyone, in our case we want only a small amount
of changes so the yaml content below should suffice for the installation in the
development environment.


```yaml
# values.yaml
# contains external-dns config

# allows to delete the dns registry when deleting a service
policy: sync 

# specifies the suffix, only will act on ingresses and services that have their
# host/annotation ending with this domain
clusterDomain: inspr.dev

# cloud provider configuration
provider: google
google: 
    project: insprlabs
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


