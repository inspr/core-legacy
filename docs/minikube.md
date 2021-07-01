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
Initialize the cluster by running `minikube start`. Run `kubectl config get-contexts` to confirm that the currently selected cluster is the minikube instance.

After all is well and done run the command below for setting up the permissions of the cluster created and to create the namespace for the inspr-apps.

```
$ kubectl create clusterrolebinding --clusterrole admin --serviceaccount default:default defaultserviceaccount
$ kubectl create namespace inspr-apps
```

## Installation of the insprd and uidp

The process of installation involves adding a few helm repositories and running `helm install`, the detailed process can be found [here](helm_installation.md).

In short the process consists of the following commands:

### **For nginx**
```
$ helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
$ helm repo update
$ helm install ingress-nginx ingress-nginx/ingress-nginx
```

Alternatively on minikube you could use the command
```
$ minikube addons enable ingress
```

For confirmation that the service is up and running run the command below, it should be a service with the name `ingress-nginx-controller`.
```
$ kubectl get service -A
```

### **For Insprd**
```bash
$ helm repo add inspr https://inspr-charts.storage.googleapis.com/
$ helm repo update
$ helm install insprd inspr/insprd \
    --set insprdRepository=gcr.io/red-inspr/insprd:v0.1.0 \
    --set authRepository=gcr.io/red-inspr/authsvc:v0.1.0 \
    --set secretRepository=gcr.io/red-inspr/secretgen:v0.1.0 \
    --set lbsidecarImage=gcr.io/red-inspr/inspr/sidecar/lbsidecar:v0.1.0
```

### **For the UIDP**

With the `insprd` and its `auth-svc` installed we can now install the inspr's UIDP. This can be done via skaffold and the [github repository](https://github.com/inspr/inspr) that contains the necessary information.


Firstly we need to open a ip so we can make requests to our minikube services, this can be done by the commands bellow.

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

Clone the github repository containing the UIDP information and alter the file `build/uidp_helm/values.yaml`, set the value of admin token to be the the one that you got from the insprctl command.

### installing redis
Install redis in the cluster via the yamls that can be found on the [github repo](https://github.com/inspr/inspr/tree/develop/cmd/uid_provider/k8s_redisdep).

Run the following commands

```bash
kubectl apply -f uidp-dep.yaml
kuebctl apply -f uidp-svc.yaml

# wait for the deployment of the 6 pods to be done
# you can check running `kubectl get pods`

# then run the following command
kubectl exec -it redis-cluster-0 -- redis-cli --cluster create --cluster-replicas 1 $(kubectl get pods -l app=redis-cluster -o jsonpath='{range.items[*]}{.status.podIP}:6379 ' | sed -e 's/ :6379/ 127.0.0.1:6379/')
```


### uidp install
Now that the admin token is set enter the upper directory of the cloned repository and run the command `skaffold run --profile uidp`.

>This will install the uidp into the minikube cluster.

With everything ready for the creation of dApps in the minikube cluster we only need to create a user in our system, that can be done by using the `inprov` cli.

> firstly create a `create_user.yaml` by using the command below
```bash
echo "
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
" > user_example.yaml
```

> afterwards run the following inprov command

```bash
inprov create --yaml user_example.yaml admin 123456
```

> This will create a user that has permissions to do any CRUD operation with dApps in the cluster. Futhermore it is a permanent account, unlike the admin initial user.
> 
> To login use the following command
```bash
inprov login minikube 123 
```


### Building docker images

For the deployment of docker images in the dapps, one has to have a public repository for docker images. The one we recommend is the dockerhub, for using the plataform follow this [tutorial](https://docs.docker.com/docker-hub/repos/#:~:text=To%20push%20an%20image%20to,docs%2Fbase%3Atesting%20)


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

In the pingpong folder there are golang files that can be built into docker images, since dApps take docker images as references we shall do the procedure to of creating those images and pushing them to ther public repository.


> creating the docker images

Clone the inspr repository, and open the terminal in the `pingpong_demo`.

Run the following commands

```bash
### builds the ping docker image and send it to the docker hub
$ docker build -t <public_repo_url>/ping:latest -f ping/ping.Dockerfile ../..
$ docker push <public_repo_url>/ping:latest

### builds the pong docker image and send it to the docker hub
$ docker build -t <public_repo_url>/pong:latest -f pong/pong.Dockerfile ../..
$ docker push <public_repo_url>/pong:latest
```

After doing the previous steps, it is now necessary to change the `image` field in `ping` and `pong` yaml files in the `pingpong_example/yamls/nodes` folder. 

> the new value should be the `dockerhub_url/<image_name>:tag`, just like the it was typed in the docker push in the steps above.


### Deploying dApps

Run the following commands to deploy dapps into your minikube cluster.

```bash
$ echo "
bootstrapServers: kafka.default.svc:9092
autoOffsetReset: earliest
sidecarImage: gcr.io/insprlabs/inspr/sidecar/kafka:latest
sidecarAddr: "http://localhost"
" > kafkaConfig.yaml

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
