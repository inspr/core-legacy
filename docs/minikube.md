# Minikube installation


For the purporse of testing Insprd locally we recommmend using [minikube](https://minikube.sigs.k8s.io/docs/start/).

Following this tutorial you will see the necessary steps to have one of the Inspr's examples running locally, more specifically the `ping pong` demonstration, but it will work for any other.

The requirements for this turorial are:
- Docker Hub account
- Docker installed
- kubectl installed
- Helm installed
- insprctl installed
- Skaffold installed


# Step by Step

## Setting up minikube
Initialize the cluster by running `minikube start` in your terminal. 

Run `kubectl config get-contexts` to confirm that the currently selected cluster is minikube.

After this is set, run the commands below for setting up the permissions of the cluster created and to create the namespace `inspr-apps`.

```
$ kubectl create clusterrolebinding --clusterrole admin --serviceaccount default:default defaultserviceaccount
$ kubectl create namespace inspr-apps
```

## Installation of the insprd and uidp

The process of installation involves adding a few Helm repositories and running `helm install`, the detailed process can be found [here](helm_installation.md).

In short the process consists of the following commands:

### **Installing NGINX**
```
$ helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
$ helm repo update
$ helm install ingress-nginx ingress-nginx/ingress-nginx
```

For confirmation that the service is up and running run the command below, it should be a service with the name `ingress-nginx-controller`.
```
$ kubectl get service -A
```

## **Installing Insprd**
```bash
$ helm repo add inspr https://inspr-charts.storage.googleapis.com/
$ helm repo update
$ helm install insprd inspr/insprd \
    --set insprdRepository=gcr.io/red-inspr/insprd:v0.1.0 \
    --set authRepository=gcr.io/red-inspr/authsvc:v0.1.0 \
    --set secretRepository=gcr.io/red-inspr/secretgen:v0.1.0 \
    --set lbsidecarImage=gcr.io/red-inspr/inspr/sidecar/lbsidecar:v0.1.0
```

## **Installing UIDP Provider**

With the `insprd` and its `auth-svc` installed we can now install the inspr's UIDP Provider. This can be done by running `skaffold` in the [github repository](https://github.com/inspr/inspr) that contains the necessary information (you should clone the repository to your local machine).


First we need to open an IP so we can make requests to our minikube services, this can be done by the commands below.

```bash
# leave it open in another terminal
$ minikube tunnel
# get the ip in which you will be doing the request
$ minikube ip
```

With the IP in hand add the entry to your local `/etc/hosts` in the following format:

```bash
127.0.0.1       localhost
::1             localhost
$MINIKUBE_IP    inspr.com
```


### **Setting up environment variables**

Set the environment variable on your local machine to fit the URL set in the `/etc/hosts`

```bash
export INPROV_PROVIDER_URL="http://inspr.com"
```

### Obtaining the auth-svc token
After setting up the IP of minikube, use `insprctl` to obtain the `adminToken` generated in its service creation.
```
$ insprctl cluster init 1234567890
```
This will return a token in a similar format to:
```
❯ insprctl cluster init 1234567890
This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjUwNjkyNzQsInBheWxvYWQiOnsidWlkIjoiIiwicGVybWlzc2lvbnMiOnsiIjpbImNyZWF0ZTp0b2tlbiJdfSwicmVmcmVzaCI6bnVsbCwicmVmcmVzaHVybCI6IiJ9fQ.IHewUCK1kGZ5OspsSr9q-yzasvYOvRbKK887PAa_iFrFkE088FkpbnyffG2z4MiuFxypGxALuZ9CVW9MtREJrQ
```
Alter the file `build/uidp_helm/values.yaml` in the repository you cloned, set the value of `adminToken` to be the one that you got from the `insprctl` command.

### **Installing Redis**
To install Redis in the cluster we need the YAML files that we can use to create the Service and Configmap.

Create two files with the content provided below.
```bash
# content of configMap yaml file named uidp-configMap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-cluster
data:
  update-node.sh: |
    #!/bin/sh
    REDIS_NODES="/data/nodes.conf"
    sed -i -e "/myself/ s/[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}/${POD_IP}/" ${REDIS_NODES}
    exec "$@"
  redis.conf: |+
    cluster-enabled yes
    cluster-require-full-coverage no
    cluster-node-timeout 15000
    cluster-config-file /data/nodes.conf
    cluster-migration-barrier 1
    appendonly yes
    protected-mode no
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis-cluster
spec:
  serviceName: redis-cluster
  replicas: 6
  selector:
    matchLabels:
      app: redis-cluster
  template:
    metadata:
      labels:
        app: redis-cluster
    spec:
      containers:
      - name: redis
        image: redis:6.2.1-alpine
        ports:
        - containerPort: 6379
          name: client
        - containerPort: 16379
          name: gossip
        command: ["/conf/update-node.sh", "redis-server", "/conf/redis.conf"]
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        volumeMounts:
        - name: conf
          mountPath: /conf
          readOnly: false
        - name: data
          mountPath: /data
          readOnly: false
      volumes:
      - name: conf
        configMap:
          name: redis-cluster
          defaultMode: 0755
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 50Mi
```

```bash
# content of the service yaml file named uidp-svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-cluster
spec:
  type: ClusterIP
  ports:
  - port: 6379
    targetPort: 6379
    name: client
  - port: 16379
    targetPort: 16379
    name: gossip
  selector:
    app: redis-cluster
```

After creating the files run the following commands:

```bash
kubectl apply -f uidp-configMap.yaml
kuebctl apply -f uidp-svc.yaml

# wait for the deployment of the 6 pods to be done
# you can check running `kubectl get pods`

# then run the following command
kubectl exec -it redis-cluster-0 -- redis-cli --cluster create --cluster-replicas 1 $(kubectl get pods -l app=redis-cluster -o jsonpath='{range.items[*]}{.status.podIP}:6379 ' | sed -e 's/ :6379/ 127.0.0.1:6379/')
```

### **Using Skaffold**
Now that the adminToken is set enter the root directory of the cloned repository and run the command `skaffold run --profile uidp`. This will install the UID Provider into the minikube cluster.

## **Deploying dApps**

### **Setting up permissions' account**

With everything ready for the creation of dApps in the minikube cluster we only need to create a user in our system, that can be done by using the `inprov` CLI. It can be installed by running the following command from within the cloned repository's root folder:

```bash
go install ./cmd/uid_provider/inprov
```

First create a `create_user.yaml` with the following content:
```bash
# content of the create_user.yaml
uid: minikube
password: "123"
permissions:
  "":
    - "create:token"

    - "create:dapp"
    - "create:channel"
    - "create:type"
    - "create:alias"

    - "get:dapp"
    - "get:channel"
    - "get:type"
    - "get:alias"

    - "update:dapp"
    - "update:channel"
    - "update:type"
    - "update:alias"

    - "delete:dapp"
    - "delete:channel"
    - "delete:type"
    - "delete:alias"
```


Run the following `inprov` command:

```bash
inprov create --yaml create_user.yaml admin 123456
```

This will create a user that has permissions to do any CRUD operation with dApps in the cluster.

To login use the following command
```bash
inprov login minikube 123 
```

### **Installing a Message Broker**

With everything in place we now need a structure in which messages between dApps can be send through, for this example we will use `Kafka`.

To install Kafka use the following commands:
```bash
$ helm repo add bitnami https://charts.bitnami.com/bitnami

$ helm install my-release bitnami/kafka
```

### **Docker Images**

For the deployment of Docker images in the dApps, one has to have a public repository for Docker images. The one we recommend is the Docker Hub, for using the platform follow this [tutorial](https://docs.docker.com/docker-hub/repos/#:~:text=To%20push%20an%20image%20to,docs%2Fbase%3Atesting%20).


The reason why this is necessary is because the dApp YAML searches for the Docker image in the repository given to it in its definition. In the pingpong example one of the definitions looks like this:
```yaml
apiVersion: v1
kind: dapp

meta:
  name: ping
  parent: pptable
spec:
  node:
    spec:
      replicas: 1
      image: gcr.io/insprlabs/inspr/example/ping:latest
      environment:
        SUPER_SECRET_0001: "false"
  boundary:
    input:
      - pinginput
    output:
      - pingoutput
```

If you like, you can use the repository that already contains the pingpong examples folder, but for better understanding of how the development process with the Inspr platform works it is recommended to follow the steps below.

In the pingpong folder are Golang files that can be built into Docker images, since dApps take those images as references we shall do the procedure of creating them and pushing it to the public repository.

### Creating the Docker Images

Open the terminal and, in the cloned Inspr repository, access the `/examples/pingpong_demo` directory.

Run the following commands
> `dockerhub_url/<image_name>:tag` is related to the repository that you are using, for example when using dockerhub it would be something like `nicholasinspr/minikube-test`. A recommendation is to use the tag `latest` when building this example, since it facilitates the process.

```bash
### builds the ping docker image and send it to the docker hub
$ docker build -t <public_repo_url>-ping:latest -f ping/ping.Dockerfile ../..
$ docker push <public_repo_url>-ping:latest

### builds the pong docker image and send it to the docker hub
$ docker build -t <public_repo_url>-pong:latest -f pong/pong.Dockerfile ../..
$ docker push <public_repo_url>-pong:latest
```

After doing the previous steps, it is now necessary to change the `image` field in `ping` and `pong` YAML files in the `pingpong_example/yamls/nodes` folder.

The new value should be the `dockerhub_url/<image_name>:tag`, just like it was typed in the docker push the steps above.
### Deploying dApps

Run the following commands to deploy dapps into your minikube cluster.

```bash
$ echo '
bootstrapServers: kafka.default.svc:9092
autoOffsetReset: earliest
sidecarImage: gcr.io/insprlabs/inspr/sidecar/kafka:latest
sidecarAddr: "http://localhost"
' > kafkaConfig.yaml

### installs message broker into the cluster
$ insprctl cluster config kafka kafkaConfig.yaml

### deploys channels' type
$ insprctl apply -k yamls/ctypes
### deploys channels
$ insprctl apply -k yamls/channels
### deploys the table
$ insprctl apply -f yamls/table.yaml
### deploys the nodes
$ insprctl apply -k yamls/nodes
```

It's important to make it clear that in the YAML files that are in the nodes, `ping.app.yaml` and `pong.app.yaml`, their field image has to be change to contain your repository where you stored the Docker image.

For illustration the `ping.app.yaml` should be like:

```yaml
apiVersion: v1
kind: dapp

meta:
  name: ping
  parent: pptable
spec:
  node:
    spec:
      replicas: 1
      # here is a example of the change that has to happen
      image: nicholasinspr/minikube-test-ping:latest
      environment:
        SUPER_SECRET_0001: "false"
  boundary:
    input:
      - pinginput
    output:
      - pingoutput

```

And that's it! You can now access your minikube cluster and see ping and pong nodes running!
