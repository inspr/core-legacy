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

####
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

#### installing redis
Install redis in the cluster via the yamls that can be found on the 



#### helm install
Now that the admin token is set enter the upper directory of the cloned repository and run the command `skaffold run --profile uidp`.
