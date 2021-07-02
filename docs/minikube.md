# Minikube installation

For the purporse of testing locally the deployment of dApps we recommmend using [minikube](https://minikube.sigs.k8s.io/docs/start/).

Following this tutorial you will see the necessary steps to have one of the inspr's examples running locally, more specifically the `ping pong` demonstration but it will work for any other.


The requirements for this turorial are:
- docker hub account
- docker installed
- kubectl installed
- helm installed
- insprctl installed
- skaffold installed


# Step by Step

## Setting up minikube
Initialize the cluster by running `minikube start`. 
>Run `kubectl config get-contexts` to confirm that the currently selected cluster is the minikube instance.

After all is well and done run the commands below for setting up the permissions of the cluster created and to create the namespace for the inspr-apps.

```
$ kubectl create clusterrolebinding --clusterrole admin --serviceaccount default:default defaultserviceaccount
$ kubectl create namespace inspr-apps
```

## Installation of the insprd and uidp

The process of installation involves adding a few helm repositories and running `helm install`, the detailed process can be found [here](helm_installation.md).

In short the process consists of the following commands:

### **Using Helm**
```
$ helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
$ helm repo update
$ helm install ingress-nginx ingress-nginx/ingress-nginx
```

Alternatively on **minikube** you could use the command
```
$ minikube addons enable ingress
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

## **Installing UIDP**

With the `insprd` and its `auth-svc` installed we can now install the inspr's UIDP. This can be done by running `skaffold` in the [github repository](https://github.com/inspr/inspr) that contains the necessary information.


Firstly we need to open an ip so we can make requests to our minikube services, this can be done by the commands below.

```bash
# leave it open in another terminal
$ minikube tunnel
# get the ip in which you will be doing the request
$ minikube ip
```

With the ip in hand add the entry to your local `/etc/hosts` in the following format:

```bash
127.0.0.1       localhost
::1             localhost
$MINIKUBE_IP    inspr.com
```

### Obtaining the auth-svc token
After setting up the ip of minikube, use the `insprctl` to obtain the adminToken generated in its service creation.
```
$ insprctl cluster init 1234567890
```
This will return a token in a similar format to:
```
❯ insprctl cluster init 1234567890
This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjUwNjkyNzQsInBheWxvYWQiOnsidWlkIjoiIiwicGVybWlzc2lvbnMiOnsiIjpbImNyZWF0ZTp0b2tlbiJdfSwicmVmcmVzaCI6bnVsbCwicmVmcmVzaHVybCI6IiJ9fQ.IHewUCK1kGZ5OspsSr9q-yzasvYOvRbKK887PAa_iFrFkE088FkpbnyffG2z4MiuFxypGxALuZ9CVW9MtREJrQ
```

Clone the github repository containing the UIDP information and alter the file `build/uidp_helm/values.yaml`, set the value of `adminToken` to be the the one that you got from the insprctl command.

### **Installing Redis**
To install redis in the cluster we need the yaml files that we can use to create the service and configMap. 

In you terminal create the two necessary files by using the commands below.
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
" > uidp-svc.yaml
```


Run the following commands

```bash
kubectl apply -f uidp-configMap.yaml
kuebctl apply -f uidp-svc.yaml

# wait for the deployment of the 6 pods to be done
# you can check running `kubectl get pods`

# then run the following command
kubectl exec -it redis-cluster-0 -- redis-cli --cluster create --cluster-replicas 1 $(kubectl get pods -l app=redis-cluster -o jsonpath='{range.items[*]}{.status.podIP}:6379 ' | sed -e 's/ :6379/ 127.0.0.1:6379/')
```


### **Using Skaffold**
Now that the admin token is set enter the upper directory of the cloned repository and run the command `skaffold run --profile uidp`.

>This will install the uidp into the minikube cluster.


## **Deploying dApps**

### **Setting up permissions' account**

With everything ready for the creation of dApps in the minikube cluster we only need to create a user in our system, that can be done by using the `inprov` cli.

> firstly create a `create_user.yaml` by using the command below
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


Run the following inprov command

```bash
inprov create --yaml create_user.yaml admin 123456
```

> This will create a user that has permissions to do any CRUD operation with dApps in the cluster. Futhermore it is a permanent account, unlike the admin initial user.
> 

To login use the following command
```bash
inprov login minikube 123 
```


### **Docker Images**

For the deployment of docker images in the dapps, one has to have a public repository for docker images. The one we recommend is the dockerhub, for using the plataform follow this [tutorial](https://docs.docker.com/docker-hub/repos/#:~:text=To%20push%20an%20image%20to,docs%2Fbase%3Atesting%20).


The reason why such this is necessary is because the dApp yaml searches for the docker image in the repository given to him in his definition, for example in the pingpong example one of the definitions looks like this:
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

If you like, you can use the repository that is already contained in the pingpong examples folder but for better understanding how the development process with the inspr plataform works it is recommended to follow the steps below.

In the pingpong folder there are golang files that can be built into docker images, since dApps takes those images as references we shall do the procedure to of creating them and pushing it to the public repository.


### Creating the Docker Images

Clone the inspr repository, and open the terminal in the `pingpong_demo` directory.

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

After doing the previous steps, it is now necessary to change the `image` field in `ping` and `pong` yaml files in the `pingpong_example/yamls/nodes` folder. 

> the new value should be the `dockerhub_url/<image_name>:tag`, just like the it was typed in the docker push in the steps above.


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
