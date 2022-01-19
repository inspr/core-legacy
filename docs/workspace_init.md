# Workspace Initialization

This document is a detailed guide on how to create an application and deploy an Inspr workspace that cointains it into a Kubernetes cluster.

**It is mandatory to have [Docker](https://docs.docker.com/get-docker/), [Kubernetes](https://kubernetes.io/docs/tasks/tools/), [Skaffold](https://skaffold.dev/docs/install/), [Insprd](helm_installation.md) and [Inspr CLI](cli-install.md) installed.**  
Also, this tutorial will be using:
- Message Broker [Apache Kafka](https://kafka.apache.org/). You can see how to install it in your cluster [here](https://bitnami.com/stack/kafka/helm).
- Redis as the database to store Insprd (Inspr daemon) users. You can install it by applying the YAML config files in [cmd/uid_provider/k8s_redisdep/](../cmd/uid_provider/k8s_redisdep/):  
		1. `kubectl apply -f cmd/uid_provider/k8s_redisdep/uidp-dep.yaml`  
		2. `kubectl apply -f cmd/uid_provider/k8s_redisdep/uidp-svc.yaml`  
		3. `kubectl exec -it redis-cluster-0 -- redis-cli --cluster create --cluster-replicas 1 $(kubectl get pods -l app=redis-cluster -o jsonpath='{range.items[*]}{.status.podIP}:6379 ' | sed -e 's/ :6379/ 127.0.0.1:6379/')`

The Inspr workspace that will be created through this guide contains two applications that communicate with each other:

- **Ping**, that writes the string "Ping!" into a Channel, and reads what's written by the other application in another Channel.
- **Pong**, that writes the string "Pong!" into a Channel, and reads what's written by the other application in another Channel.

The tree-like structure of this Inspr workspace will be similar to this:

![structure](img/pingpong-app.jpg)

## Creating the folder structure

First of all, let's create the folders in which we will store all of our files (application codes, Dockerfiles, YAMLs, etc.).

Create a folder called "pingpong_demo", and inside of it create other three folders: "ping", "pong" and "yamls".

The first two, "ping" and "pong", will store the applications code and their respective Dockerfiles. The last one, "yamls", will store all the `.yaml` files that define the structures which will be used to create the Inspr workspace and structures in the cluster.

The following command creates the wanted folder structure:
```zsh
mkdir -p pingpong_demo/ping pingpong_demo/pong pingpong_demo/yamls
```

It should be organized like this:

```
pingpong_demo
├── ping
├── pong
└── yamls
```

## Creating the applications

In this part, we will implement Ping and Pong using Golang. Also, we will create their respective Dockerfiles, build the Docker Images and push them into a Docker Registry.  
To start off, inside of "/pingpong_demo", run the following command:
```go
go mod init pingpong
```

If you're not familiar with `go mod`, you can learn more about it [here](https://golang.org/ref/mod#go-mod-init).

### Ping and Pong implementation

From within "pingpong_demo", create a file called `ping.go` inside the folder "ping":
```zsh
touch ping/ping.go
```

In `ping.go`, we will define a `main` function that does the following:

1. Creates a new dApp Client, which is used to write and read messages in Channels through the Sidecar (check [dApp Architecture Overview](dapp_overview.md) for more details).
2. Call the Client's `WriteMessage` method to write the string "Ping!" in the Channel _pingoutput_.
3. Define the Client's handler, which will read messages from the Channel _pinginput_ and print it. Then it proceeds to write a new message just like described in the previous step.

The first thing we want to do in `main.go` is to create a new dApp Client to enable the Node-Sidecar communication, and set the message we want to send:
```go
func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentMsg := "Ping!"
}
```

Now we must send the first message to _pingoutput_, and handle any error that may occur:
```go
if err := client.WriteMessage(ctx, "pingoutput", sentMsg); err != nil {
	fmt.Printf("an error occurred: %v", err)
	return
}
```

We then implement the Client's handler, which will handle every message that is received through the Channel _pinginput_. It receives a message and decodes it into a structure called `BrokerMessage`. This `BrokerMessage` contains only one field called `Data`, in which our received message will be stored. Then it proceeds to print the received message, and sends a new message just like in the previous step:
```go
client.HandleChannel("pinginput", func(ctx context.Context, body io.Reader) error {
	var ret models.BrokerMessage

	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&ret); err != nil {
		return err
	}

	fmt.Println(ret)

	if err := client.WriteMessage(ctx, "pingoutput", sentMsg); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
})
```

And finally, we start our Client so it's able to receive messages:
```go 
log.Fatal(client.Run(ctx))
```

When completed, your `ping.go` implementation should look like this:
```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
	"golang.org/x/net/context"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentMsg := "Ping!"
	if err := client.WriteMessage(ctx, "pingoutput", sentMsg); err != nil {
		fmt.Printf("an error occurred: %v", err)
		return
	}
	client.HandleChannel("pinginput", func(ctx context.Context, body io.Reader) error {
		var ret models.BrokerMessage

		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&ret); err != nil {
			return err
		}

		fmt.Println(ret)

		if err := client.WriteMessage(ctx, "pingoutput", sentMsg); err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	log.Fatal(client.Run(ctx))
}
```

A similar folder/file structure and code must be done to implement Pong. So, from within "pingpong_demo", create `pong.go` in /pong folder:
```zsh
touch pong/pong.go
```

And then write a code similar to _ping.go_'s, just remember to switch the Channels and the message that is written. It should look like this:
```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
	"golang.org/x/net/context"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sentMsg := "Pong!"
	if err := client.WriteMessage(ctx, "pongoutput", sentMsg); err != nil {
		fmt.Printf("an error occurred: %v", err)
		return
	}
	client.HandleChannel("ponginput", func(ctx context.Context, body io.Reader) error {
		var ret models.BrokerMessage

		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&ret); err != nil {
			return err
		}

		fmt.Println(ret)

		if err := client.WriteMessage(ctx, "pongoutput", sentMsg); err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	log.Fatal(client.Run(ctx))
}
```

To finish implementing the applications, we must assure that all dependencies are resolved. To do so, inside of "/pingpong_demo", run the following command:
```go
go mod tidy
```

### Dockerfiles

Now that we have Ping and Pong implemented and doing what they are supposed to do, we must create their Dockerfiles so a Docker Image for each of them can be generated, and make these images available in a registry such as [Docker Hub](https://hub.docker.com/).  
If you are not familiar with Docker, click on the following links for more information:

- [Docker overview](https://docs.docker.com/get-started/overview/)
- [Docker Images](https://jfrog.com/knowledge-base/a-beginners-guide-to-understanding-and-building-docker-images/#:~:text=A%20Docker%20image%20is%20a,publicly%20with%20other%20Docker%20users.)
- [Dockerfiles](https://docs.docker.com/engine/reference/builder/)

From within "pingpong_demo", create a file called `ping.Dockerfile` inside folder "ping":
```zsh
touch ping/ping.Dockerfile
```

The Dockerfile structure will be created to do the following:

1. Use Golang alpine version as base to run the commands.
2. Define "/app" as the working directory where Ping will be build.
3. Copy the `ping.go` file into "/app".
4. Build `ping.go` to generate a binary for this file.
5. Execute the binary generated.

Ping's Dockerfile should look like this:
```docker
FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go build -o main pingpong_demo/ping/ping.go

FROM alpine
WORKDIR /app
COPY --from=build /app/main .
CMD ./main
```

Then, do the same steps for Pong:
```zsh
touch pong/pong.Dockerfile
```

And Pong's Dockerfile content:
```docker
FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go build -o main pingpong_demo/pong/pong.go

FROM alpine
WORKDIR /app
COPY --from=build /app/main .
CMD ./main
```

### Docker Image deployment

After creating the Dockerfiles, we must build the Docker images and push each of them into the Docker Registry, so they are available to be used in the Inspr workspace.

In the next steps, it's important to have a [container registry](https://cloud.google.com/container-registry) up and running, so it's possible to store and use the images that will be built.

Then, from the "pingpong_demo" folder, we must build the Docker image by using the Dockerfile previously created. You can apply a tag to it by adding `:TAG_NAME`, if desired:
```zsh
docker build -f ping/Dockerfile -t CONTAINER_REGISTRY_REF/app/ping:TAG_NAME .
```

Finally, we push the builded image into the registry, where **CONTAINER_REGISTRY_REF is a reference to the container registry**:
```zsh
docker push CONTAINER_REGISTRY_REF/app/ping:TAG_NAME
```

Now the same must be done to Pong:
```zsh
docker build -f pong/Dockerfile -t CONTAINER_REGISTRY_REF/app/pong:TAG_NAME .
```

Push the builded image into the registry:
```zsh
docker push CONTAINER_REGISTRY_REF/app/pong:TAG_NAME
```

**Alternatively, you can create a [Makefile](https://opensource.com/article/18/8/what-how-makefile) that will do the previous steps for you.**  
Inside "/pingpong_demo" folder, create a new file called "Makefile:
```zsh
touch Makefile
```

The Makefile should contain the same Docker commands that you'd use to build and push Ping and Pong Docker Images "manually". The gain here is that instead of writing and executing four different commands, you just execute the Makefile. It should look like this:
```makefile
build:
	docker build -t CONTAINER_REGISTRY_REF/app/pong:TAG_NAME -f ping/Dockerfile .
	docker push CONTAINER_REGISTRY_REF/app/pong:TAG_NAME
	docker build -t CONTAINER_REGISTRY_REF/app/pong:TAG_NAME -f pong/Dockerfile .
	docker push CONTAINER_REGISTRY_REF/app/pong:TAG_NAME
```

And to execute the Makefile through the terminal you just run:
```zsh
make
```

## Creating Inspr Workspace

Now that we have our applications implemented and their Docker images available in the cluster, we're good to build the Inspr structures to run the applications we created.

### YAML Files

First of all, from "/pingpong_demo" we access the folder "/yamls" created previously.
```zsh
cd yamls
```

And within this folder we will create the `.yaml` files which describe each of the Inspr structures that will be built inside the cluster. This part of the tutorial won't take a closer look at every minimum detail on how to write the YAML files, but you can find more information about it [here](yamls/index.md).

**1) dApp YAMLs**  
The first file to be created is `table.yaml`, which is the dApp that will contain Ping and Pong Nodes:
```zsh
touch table.yaml
```

As it's described in Inspr YAMLs documentation, we must specify the kind, apiVersion and then the dApp information. This specific dApp will work as a link between Ping/Pong and the Channels which they'll communicate through. Basically, it will connect Ping and Pong's Boundaries to Channels defined in the root dApp through Aliases. You can read more about Aliases and this dApp structure [here](dapp_overview.md).  
The YAML should be the following:
```yaml
apiVersion: v1
kind: dapp

meta:
  name: pptable
spec:
  aliases:
    ping.pingoutput:
      resource: ppchannel1
    ping.pinginput:
      resource: ppchannel2
    pong.ponginput:
      resource: ppchannel1
    pong.pongoutput:
      resource: ppchannel2
  boundary:
    input:
      - ppchannel2
      - ppchannel1
    output:
      - ppchannel2
      - ppchannel1
```

The next file to be created is _ping.app.yaml_, in which we describe the Node that contains Ping. For better organization, let's create a new folder inside of "/yamls" called "nodes", inside of which we will store the Nodes YAMLs:
```zsh
mkdir nodes
touch nodes/ping.app.yaml
```

Then, inside of `ping.app.yaml` we must specify the kind, apiVersion and then the Node information (such as name, boundaries, image, etc.).  
It should look like this:
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
      image: CONTAINER_REGISTRY_REF/app/ping:TAG_NAME
	  environment:
        SUPER_SECRET_0001: "false"
  boundary:
    input:
      - pinginput
    output:
      - pingoutput
```

The `image` field must be a reference to the same image built in the step **"Docker Image deployment"**. Also, the `environment` field is defined just so we can see that the environment variable "SUPER_SECRET_0001" will be created in the cluster, inside Ping's deployment.

Now, we do the same for Pong:
```zsh
touch nodes/pong.app.yaml
```

And `pong.app.yaml` should look like this:
```yaml
apiVersion: v1
kind: dapp

meta:
  name: pong
  parent: pptable
spec:
  node:
    spec:
      image: CONTAINER_REGISTRY_REF/app/pong:TAG_NAME
  boundary:
    input:
      - ponginput
    output:
      - pongoutput
```

Notice that we didn't specify the number of replicas for Pong. This is to show that, **if not specified, the default number of replicas is created, which is one replica.**

**2) Channel YAMLs**  
As it can be seen in dApp `table.yaml`, the four different Boundaries defined in Ping and Pong (`pinginput`, `pingoutput`, `ponginput` and `pongoutput`) actually refer only to two Channels, `ppchannel1` and `ppchannel2`, which need to create in the root dApp.  
Similar to a dApp, we must specify the kind, apiVersion and then the Channel information. Notice that the Channel `brokerlist` has only Kafka listed, as we're using only Kafka as our Message Broker.

Let's proceed to create "/channels" folder inside of "/yamls", and then create `ch1.yaml`:
```zsh
mkdir channels
touch channels/ch1.yaml
```

It's content should be:
```yaml
apiVersion: v1
kind: channel

meta:
  name: ppchannel1
spec:
  type: pptype1
  brokerlist:
    - kafka
```

Now we do the same for `ch2.yaml`:
```zsh
touch channels/ch2.yaml
```

And it's content should be:
```yaml
apiVersion: v1
kind: channel

meta:
  name: ppchannel2
spec:
  type: pptype1
  brokerlist:
    - kafka
```

**3) Type YAML**  
Both Channels use the same Type to define the kind of message that goes through them. So we must specify the kind, apiVersion and then the Type information for it to be created in the same context (the same dApp) as the Channels'.  
First we create "/types" folder and `ct1.yaml`:
```zsh
mkdir types
touch types/ct1.yaml
```

And it's content should be:
```yaml
apiVersion: v1
kind: type

meta:
  name: pptype1
schema: yamls/ctypes/schema.avsc
```

Notice that the `schema` field is actually a reference to an **Avro Schema** file. By defining it like this, when a Type is created Inspr searches for the file and injects its value into the `schema` field. You can find more information on how schemas should be created to be used in Inspr [here](schemas_and_types.md).  
To make everything work properly, let's create `schema.avsc`:
```zsh
touch types/schema.avsc
```

As we defined in our Ping and Pong applications, the type of information that they will send and receive is just a simple string. To do so, `schema.avsc` content should be:
```
{"type":"string"}
```

**4) Kafka Sidecar Configuration YAML**  
As our Channels will communicate through Kafka Message Broker, the Kafka's Sidecar configurations must be defined.  
In the "yamls/" folder create the file `kafka.yaml`:
```zsh
touch kafka.yaml
```

And in this file we will set some basic configs for our Kafka Sidecar:
- Bootstrap Server: basically an in-cluster reference to your Kafka Message Broker service
- Auto Offset Reset: see more in [Kafka's documentation](https://kafka.apache.org/documentation.html#consumerconfigs_auto.offset.reset)
- Sidecar Image: references the Docker Image of the Sidecar to be used (you probably wont need to change this)
- Sidecar Address: the address in the k8s Pod in which the Sidecar will be (it's a good idea to keep it as localhost)  

So it will look like this:
```yaml
bootstrapServers: <kafka service name>.default.svc:9092
autoOffsetReset: earliest
sidecarImage: gcr.io/insprlabs/inspr/sidecar/kafka:latest
sidecarAddr: "http://localhost"
```

After creating all of these folders and files, you should have this:
```
pingpong_demo
├── Makefile
├── ping
│   ├── ping.Dockerfile
│   └── ping.go
├── pong
│   ├── pong.Dockerfile
│   └── pong.go
├── README.md
└── yamls
    ├── channels
    │   ├── ch1.yaml
    │   └── ch2.yaml
    ├── ctypes
    │   ├── ct1.yaml
    │   └── schema.avsc
    ├── kafka.yaml
    ├── nodes
    │   ├── ping.app.yaml
    │   └── pong.app.yaml
    └── table.yaml
```
Finally, now that we have Ping and Pong images in the cluster and all Inspr workspace structures well-defined in YAML files, we can deploy everything that we created and see it working in our cluster.

### Connecting Insprd with the UID Provider

First of all, we need to check if **Inspr CLI** is referring to the [cluster ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) address, so we're able to send requests to it. To do so, run the command:
```zsh
insprctl config list
```

And something similar to the following should be shown:
```zsh
Available configurations:
- scope: ""
- serverip: "http://localhost:8080"
```

If the `serverip` is not your cluster ingress host you must change it:
```zsh
insprctl config serverip "CLUSTER_INGRESS_HOST"
```

And this will be printed in the terminal:
```zsh
Success: insprctl config [serverip] changed to 'CLUSTER_INGRESS_HOST'
```

When running `insprctl config list`, if an error like the following shows up, run the command `insprctl init` and then try again:
```
Invalid config file! Did you run insprctl init?
Error: open /home/<user>/.inspr/config: no such file or directory
```

Once this is done we need to install our **UID Provider**, so we're able to create users that have access to the workspace and can manipulate it's structures.  

### UID Provider installation

First of all, we must retrieve the **root token for authentication**. This is the [JWT](https://jwt.io/) token that will be used to establish the first connection between Insprd and the UID Provider, generating an Admin user which will be able to create new users, as well as to manipulate Inspr's structures inside Insprd.

If when installing Inspr's Helm Chart you didn't change the variable `deployment.initKey`, it's default value is "1234567890" (if you did change, use the value you inserted). Run the following command replacing `<key>` with the `initKey`'s value:  
```zsh
insprctl cluster init <key>
```

After that, a message similar to the following will be prompted in your terminal:
```
This is a root token for authentication within your insprd. This will not be generated again. Save it wisely.
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjM5NDk5NzMsInBheWxvYWQiOnsidWlkIjoiIiwicGVybWlzc2lvbnMiOnsiIjpbImNyZWF0ZTp0b2tlbiJdfSwicmVmcmVzaCI6bnVsbCwicmVmcmVzaHVybCI6IiJ9fQ.jDxtf_S3QLWPilwparQ3rl3skdDfR9JYKwILc08X19ZTiLA8bs9eaC9AC3o0pYI7GrmVxoei76vFWYlwCfviCA
```

To install the UID Provider in the cluster, we'll want to use **Skaffold** from within the Inspr repository. This will be made simpler in the future, but as it is now this is the only way to do so.  
So, if you haven't cloned [Inspr's repository](https://github.com/inspr/inspr) to your local machine yet, clone it. You can see how to do it [here](https://git-scm.com/book/en/v2/Git-Basics-Getting-a-Git-Repository).  

Once you have Inspr's repository cloned, go into its root folder and open the UID Provider Helm Chart `values.yaml`, it can be found in `build/uidp_helm/`. In this file, replace whatever is in `secret.adminToken` by the token you generated in the previous step. You can also alter `secret.adminPassword` if desired.  
After updating the admin token, save your changes and run the following command from within the root of the repository to install the UID Provider:
```zsh
skaffold run --profile uidp
```

When the installation is done you must export an environment variable that the UID Provider will use to reference Insprd in the cluster. This environment variable is the access address to Insprd in your cluster:
```zsh
export INPROV_PROVIDER_URL="http://<INSPRD_CLUSTER_ADDRESS>"
```

To start using the UID Provider CLI to login and manage users, you need to install it by running the following command from within the repository root:
```go
go install ./cmd/uid_provider/inprov
```

And finally, login with the admin user! The default password is "123456" if you haven't changed it in Helm Chart's `values.yaml`. 
```zsh
inprov login admin <password>
```

### Deploying dApps, Channels and Type

First let's install Kafka in Insprd using the configs written in `kafka.yaml`, so our Channels will be able to communicate through the Kafka Broker. To do so, run the following command from within the "/pingpong_demo" folder:
```
insprctl brokers kafka yamls/kafka.yaml
```
The expected response is: `successfully installed broker on insprd`.  
Now we apply the YAML files by using Inspr CLIs commands. The structures should be applied in the following order:

1. Type `pptype1`
2. Channels `ppchannel1` and `ppchannel2`
3. dApp `pptable`
4. Nodes `ping` and `pong`

You can do so by running the following commands from within "/pingpong_demo" folder:

```
insprctl apply -k yamls/ctypes
insprctl apply -k yamls/channels
insprctl apply -f yamls/table.yaml
insprctl apply -k yamls/nodes
```

To learn more about Inspr CLI, check [this documentation.](cli/inspr.md)

If everything worked fine, Insprd will have printed a changelog similar to the following for each command written in your terminal:

```
➜  pingpong_demo ✗ insprctl apply -k yamls/ctypes 
On: 
Field                 | From       | To
Spec.Types[pptype1]   | <nil>      | {...}

Applied:
ct1.yaml | type | v1

➜  pingpong_demo ✗ insprctl apply -k yamls/channels 
On: 
Field                       | From       | To
Spec.Channels[ppchannel1]   | <nil>      | {...}
On: 
Field                       | From       | To
Spec.Channels[ppchannel2]   | <nil>      | {...}

Applied:
ch1.yaml | channel | v1
ch2.yaml | channel | v1

➜  pingpong_demo ✗ insprctl apply -f yamls/table.yaml                       
On: pptable
Field                           | From       | To
Meta.Name                       |            | pptable
Spec.Boundary.Input             | <nil>      | ppchannel2
Spec.Boundary.Input             | <nil>      | ppchannel1
Spec.Boundary.Output            | <nil>      | ppchannel2
Spec.Boundary.Output            | <nil>      | ppchannel1
Spec.Aliases[ping.pinginput]    | <nil>      | ppchannel2
Spec.Aliases[ping.pingoutput]   | <nil>      | ppchannel1
Spec.Aliases[pong.ponginput]    | <nil>      | ppchannel1
Spec.Aliases[pong.pongoutput]   | <nil>      | ppchannel2
On: 
Field                | From       | To
Spec.Apps[pptable]   | <nil>      | {...}

Applied:
yamls/table.yaml | dapp | v1

➜  pingpong_demo ✗ insprctl apply -k yamls/nodes     
On: pptable.ping
Field                                           | From       | To
Meta.Name                                       |            | ping
Meta.Parent                                     |            | pptable
Spec.Node.Meta.Name                             |            | ping
Spec.Node.Meta.Parent                           |            | pptable
Spec.Node.Spec.Image                            |            | gcr.io/insprlabs/inspr/example/ping:latest
Spec.Node.Spec.Replicas                         | 0          | 1
Spec.Node.Spec.Environment[SUPER_SECRET_0001]   | <nil>      | false
Spec.Boundary.Input                             | <nil>      | pinginput
Spec.Boundary.Output                            | <nil>      | pingoutput
On: pptable
Field             | From       | To
Spec.Apps[ping]   | <nil>      | {...}
On: pptable.pong
Field                   | From       | To
Meta.Name               |            | pong
Meta.Parent             |            | pptable
Spec.Node.Meta.Name     |            | pong
Spec.Node.Meta.Parent   |            | pptable
Spec.Node.Spec.Image    |            | gcr.io/insprlabs/inspr/example/pong:latest
Spec.Boundary.Input     | <nil>      | ponginput
Spec.Boundary.Output    | <nil>      | pongoutput
On: pptable
Field             | From       | To
Spec.Apps[pong]   | <nil>      | {...}

Applied:
ping.app.yaml | dapp | v1
pong.app.yaml | dapp | v1
```

## And, finally, it's done! You have just initialized your Inspr workspace!

To check if everything worked and that Ping and Pong deployments were created, run the following command in your terminal:

```zsh
kubectl get deploy --namespace inspr-apps
```

And a message like this should be displayed:

```
NAME           								READY   UP-TO-DATE   AVAILABLE   AGE
node-5d87b3b9-611c-49e5-b76e-4838abd363f9   1/1     1            1           12m
node-c77cc79d-756c-4a88-a827-579725d8a1fc   1/1     1            1           12m
```

To check if Ping and Pong are actually writing and reading messages, you can access you cluster informations by using [k9s](https://github.com/derailed/k9s).  
In k9s, access the deployments with the same name as displayed by `kubectl get deploy` (inside the namespace `inspr-apps`), then access the pod within it. The pod should contain three containers, one for the Node, one for the Load Balancer Sidecar and the other for Kafka's Sidecar:

![structure](img/k9s-container.jpg)

And by accessing the Node container, you'll be able to see the application running!

![structure](img/k9s-log.jpg)


## If you have any questions or suggestions about this tutorial feel free to join our Discord channel and talk with the team!