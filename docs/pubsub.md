# Publish / Subscribe

Pub/Sub is an asynchronous messaging service that decouples services that produce events from services that process events. It is an architectural pattern where messages are sent in a 1:n manner

Very popular and useful, Pub/Sub offers reliable message storage and real-time delivery. Inspr is the perfect environment for these applications, seen that they benefit heavily from microservice structure as well as the shared communication Channels, both part of Inspr's main features. Being able to easily publish and change the dApp's structure also allows for great expandability, especially when it comes to adding new subscriber to the service.

Now we'll take a step by step look into how to build a Pub/Sub app on Inspr:

- Laying out the foundations.
- Publisher
- Subscriber
- Connection
- Publishing it

## Foundations

This section will cover the basic structure for the application. Divided in layers, what we are building will look something like this:

![Pub/Sub](./img/pubsub.jpg)

1. Publish:
   This layer is responsible for receiving messages from the user, these will be sent to the final clients. This is achieved by the endpoint `/publish`, served by a very simple REST API. Messages received via the endpoint are configured and sent, using the Channel we will be creating for the application, to the subscribers through the broker.

2. Broker:
   This layer is handled for us by inspr! Basically it receives messages on a Channel, stores them in a queue, and makes a copy of each message available for every reader of said Channel.

3. Subscribe:
   This layer is arguably the simplest one, it hosts an arbitrary number of subscriber dApps. Each of these dApps reads from the application's Channel, receiving the messages published by the user. Once it does get a message the subscriber configures it on it's client's desired format. Each client service expects a different set of information and flags, besides the main message text, it is the subscriber's job to correctly assemble a requisition and send it to it's client.

#### On Inspr:

Now that we have an understanding of the application we'll be developing, let's see how Inspr can help us build it!

The application as a hole is only complete if you consider all the previously explained Channels. However, as was said, we don't have to worry about the message broker, Inspr takes care of that. That leaves us with publishers and subscribers. The first thing we have to do is define a scope and Channel for our application, so let's configure a dApp.

Scope app (01.app.yaml) :

```yaml
kind: dapp
apiVersion: v1

meta:
  name: pubsub
```

Note that this dApp isn't a [Node](dapp_overview.md), that's because it serves as a workspace definition, it doesn't run any code.

Next we have to create a Channel the publisher can write text to and the subscribers can read from. It's also necessary to create a Type for the Channel, so that the kind of information that goes through the Channel is well-defined. Doing so is as simple as defining the following:

Type (02.ct.yaml) :

```yaml
kind: type
apiVersion: v1

meta:
  name: pubsubct
  parent: pubsub
schema: '{"type":"string"}'
```

Channel (03.ch.yaml) :

```yaml
kind: channel
apiVersion: v1

meta:
  name: pubsubch
  parent: pubsub
spec:
  type: pubsubct
  brokerlist:
    - kafka
```

As you can see in the Channel's YAML definition, this PubSub example will be using Apache Kafka as it's Message Broker, so be sure you have [Kafka installed in your cluster](https://bitnami.com/stack/kafka/helm). That said, we must create a YAML file that configures Kafka to be used by Inspr:  
Kafka (kafka.yaml) :
```yaml
bootstrapServers: kafka.default.svc:9092
autoOffsetReset: earliest
sidecarImage: gcr.io/insprlabs/inspr/sidecar/kafka:latest
sidecarAddr: "http://localhost"
```

Then we proceed to create the publisher dApp, which will be the one receiving the requests that contains the messages to be passed along to the subscriber dApps so they can send them. The biggest difference between this dApp and the subscriber ones is that it contains the ports through which it's going to receive the requests.

```yaml
kind: dapp
apiVersion: v1

meta:
  name: pubsubapi
  parent: pubsub # parent app name
spec:
  node:
    spec:
      replicas: 1
      image: <reference to publisher image>
      ports:
        - port: 80
          targetPort: 8080
  boundary:
    output:
      - pubsubch #Channel created above
```

Now all there is left to is creating the subscribers' dApps. Since we could build any number of subscribers, depending on what clients we would like interact with, consider the following example a template for their files.

Subscribers (04-06.app.yaml) :
```yaml
kind: dapp
apiVersion: v1

meta:
  name: <subscribername>
  parent: pubsub # parent app name
spec:
  node:
    spec:
      replicas: 1
      image: <reference to this app's image>
  boundary:
    input:
      - pubsubch #Channel created above
```

These are all the basic building blocks we'll need !

## Publisher

Implementation of the publisher goes as you would expect any simple REST API to go. First we build the API's server that will support our endpoint. Following that the data model that will be received in a requisition must be defined. In this example the API expects a JSON with only one field, `message`,  which is a string. The handler is the last part, all it should do is receive the message from the user and send it to the subscribers using the Channel we created earlier.

API (api/main.go) : Entry point for the server, must be a `package main` so that it executes on the container.

```go
package main

import (
    controller "inspr.dev/inspr/examples/pubsub/api/controller"
)

var server controller.Server

// main is the server start up function
func main() {
    server.Init()
    server.Run(":8080")
}
```

Run method (api/controller/base.go):

```go
   // Run starts the server on the port given in addr
func (s *Server) Run(addr string) { // this is called by the main()
    fmt.Printf("pubsub api is up! Listening on port: %s\n", addr)
    log.Fatal(http.ListenAndServe(addr, s.Mux))
}
```

Init method (api/controller/base.go) : This is responsible for allocating the server multiplexer (mux) and inserting our handler on it.

```go
   // Init - configures the server
func (s *Server) Init() {
	s.Mux = http.NewServeMux()
	client := dappclient.NewAppClient()
	s.Mux.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		data := message{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		pubMsg := data.Message
		if err := client.WriteMessage(ctx, pubsubChannel, pubMsg); err != nil {
			fmt.Println(err)
			rest.ERROR(w, err)
		}

		rest.JSON(w, http.StatusOK, "Message sent!")
	})
}
```

## Subscribers

As most applications that run on Inspr the subscribers only execute their tasks when they receive a message through it's dApp client. Once the message has been sent to the clients the message is committed, acknowledging to the message broker that that application is done with that message and can go on to the next one.

These examples were implemented using webhooks that receive JSON objects with specific fields.

### Discord Sub (discord.main) :

```go
type DiscordMessage struct {
    Content   string `json:"content"`
    Username  string `json:"username"`
    AvatarURL string `json:"avatar_url"`
    TTS       bool   `json:"tts"`
    File      []byte `json:"file"`
    Embedded   []byte `json:"embeds"`
}

var webhook = <your webhook link>
var channel = "pubsubch" // Channel created on parent dApp

func main() {
	c := &http.Client{}
	client := dappclient.NewAppClient()
	client.HandleChannel(channel, func(ctx context.Context, body io.Reader) error {
		decoder := json.NewDecoder(body)

		subMsg := models.BrokerMessage{}
		err := decoder.Decode(&subMsg)
		if err != nil {
			return err
		}

		msg := discordMessage{
			Content:   fmt.Sprintf("%v", subMsg.Data),
			Username:  "Notifications",
			AvatarURL: "",
			TTS:       true,
		}

		msgBuff, _ := json.Marshal(msg)

		req, _ := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(msgBuff))
		head := http.Header{}
		head.Add("Content-type", "application/json")
		req.Header = head
		_, err = c.Do(req)
		if err != nil {
			return err
		}
		return nil

	})
	log.Fatalln(client.Run(context.Background()))
}
```

## Connection

Hosting a server with an entry point on your cluster doesn't make it so that you can access it remotely, to do so you must create a Kubernetes ingress point and associate it with the publisher's dApp service. Luckily this is easy to make and there is an example right here! To connect your application to these entrypoints you could remember that on the cluster an app's name is the full path to the app, separated by hyphens.

Ingress:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress

metadata:
  name: pubsub-ingress
  namespace: inspr-apps
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/issuer: "letsencrypt-prod"
spec:
  rules:
    - host: <cluster's host>
      http:
        paths:
          - path: /publish
            pathType: Prefix
            backend:
              service:
                name: <publisher dApp service name>
                port:
                  number: 80
```

Remember to replace "<publisher dApp service name>" by the name of the service which is generated when a dApp is created in the cluster.

## Publishing it

Your app is now done, but before publishing it first you must build it. Inspr requires code to be built into an accessible Docker image in order for it to be deployed on the cluster. There is an example for building on the [pingpong example](../examples/pingpong_demo/README.md). You should create simple Dockerfiles to build your dApps and push them to any regs where they are accessible your cluster. That being done is finally time to publish your Pub/Sub application. Publishing it to your insprctl cluster is easy, but first make sure your files look something like this:

```tree
pubsub
├── api
│   ├── controller
│   │   └── base.go
│   ├── Dockerfile
│   └── main.go
├── discord
│   ├── Dockerfile
│   └── main.go
├── k8s
│   └── ingress.yaml
├── Makefile
├── slack
│   ├── Dockerfile
│   └── main.go
└── yamls
    ├── 01.app.yaml
    ├── 02.ct.yaml
    ├── 03.ch.yaml
    ├── kafka.yaml
    └── nodes
        ├── 04.app.yaml
        ├── 05.app.yaml
        └── 06.app.yaml
```

Keep in mind Slack an Discord are only examples.

Deployment:
You have to apply every YAML we created for dApps, Channels, Type and Kafka by running the following commands (be sure to have the latest version of [Inspr CLI](cli_install.md) installed):
```
insprctl cluster config kafka yamls/kafka.yaml
insprctl apply -f yamls/01.app.yaml
insprctl apply -f yamls/02.ct.yaml
insprctl apply -f yamls/03.ch.yaml
insprctl apply -k yamls/nodes
```

Your ingress have to be deployed as well, run:
```
kubectl apply -f k8s/ingress.yaml
```