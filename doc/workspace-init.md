# Workspace Initialization
This document is a detailed guide on how to create an application and deploy an Inspr environment that cointains it into a Kubernetes cluster.  

**It is mandatory to have [Docker](https://docs.docker.com/get-docker/), [Kubernetes](https://kubernetes.io/docs/tasks/tools/), [Inspr](helm-install.md) and [Inspr CLI](cli-install.md) installed. [Skaffold](https://skaffold.dev/docs/install/) can also be very useful.**

The Inspr environment that will be created through this guide contains two applications that communicate with each other:
- **Ping**, that writes the string "Ping!" into a Channel, and reads what's written by the other application in another Channel.
- **Pong**, that writes the string "Pong!" into a Channel, and reads what's written by the other application in another Channel.    

The tree-like structure of this Inspr environment will be similar to this:  

![structure](img/pingpong-app.jpg)

## Creating the folder structure

First of all, let's create the folders in which we will store all of our files (application codes, Dockerfiles, YAMLs, etc.).

Create a folder called "pingpong_demo", and inside of it create other three folders: "ping", "pong" and "yamls".  

The first two, "ping" and "pong", will store the applications code and their respective Dockerfiles. The last one, "yamls", will store all the *.yaml* files that define the structures which will be used to create the Inspr environment in the cluster.  

The following command creates the wanted folder structure:
```
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

In this part, we will implement Ping and Pong using Golang. Also, we will create their respective Dockerfiles, build the Docker Images and push them into the cluster.  

### Ping and Pong implementation  

From within "pingpong_demo", create a file called *ping.go* inside folder "ping":  
```
touch ping/ping.go
```  

In *ping.go*, we will define a `main` function that does the following:

1) Creates a new dApp Client, which is used to write and read messages in Channels through the Sidecar (check [dApp Architecture Overview](dapp-overview.md) for more details).
2) Initiates an endless `for loop` in which the message "Ping!" is written in the Channel *ppChannel1*, then the application proceeds to read a message from Channel *ppChannel2*. If there is a message, it's read and displayed in the terminal.

*ping.go* should look like this:
```
package main

import (
	"fmt"
	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	"golang.org/x/net/context"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		sentMsg := models.Message{
			Data: "Ping!",
		}

		if err := client.WriteMessage(ctx, "ppChannel1", sentMsg); err != nil {
			fmt.Println(err)
			continue
		}

		recMsg, err := client.ReadMessage(ctx, "ppChannel2")
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Read message: ")
		fmt.Println(recMsg.Data)

		if err := client.CommitMessage(ctx, "ppChannel2"); err != nil {
			fmt.Println(err.Error())
		}
	}
}
```

A similar folder/file structure and code must be done to implement Pong. So, from within "pingpong_demo", create *pong.go* in /pong folder:  
```
touch pong/pong.go
```  

And then write a code similar to *ping.go*'s, just remember to swich the Channels and the message that is written. It should look like this:
```
package main

import (
	"fmt"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
	"golang.org/x/net/context"
)

func main() {

	client := dappclient.NewAppClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		sentMsg := models.Message{
			Data: "Pong!",
		}

		if err := client.WriteMessage(ctx, "ppChannel2", sentMsg); err != nil {
			fmt.Println(err)
			continue
		}

		recMsg, err := client.ReadMessage(ctx, "ppChannel1")
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Read message: ")
		fmt.Println(recMsg.Data)

		if err := client.CommitMessage(ctx, "ppChannel1"); err != nil {
			fmt.Println(err.Error())
		}
	}
}
```

### Dockerfiles
Now that we have Ping and Pong implemented and doing what they are supossed to do, we must create their Dockerfiles, so we can generate a Docker Image for each of them, and make these images available in our cluster.  
If you are not familiar with Docker, click on the following for more information:
- [Docker overview](https://docs.docker.com/get-started/overview/)
- [Docker Images](https://jfrog.com/knowledge-base/a-beginners-guide-to-understanding-and-building-docker-images/#:~:text=A%20Docker%20image%20is%20a,publicly%20with%20other%20Docker%20users.)
- [Dockerfiles](https://docs.docker.com/engine/reference/builder/)  

From within "pingpong_demo", create a file called *Dockerfile* inside folder "ping":  
```
touch ping/Dockerfile
```  
The Dockerfile structure will be created to do the following:
1) Use Golang alpine version as base to run the commands.
2) Define "/app" as the working directory where Ping will be build.
3) Copy the *ping.go* file into "/app".
4) Compile and build the *ping.go* file.

Ping's Dockerfile should look like this:
```
FROM golang:alpine
WORKDIR /app
COPY . .
CMD go run pingpong_demo/ping/ping.go
```

Then, do the same steps for Pong:
```
touch pong/Dockerfile
```  

And Pong's Dockerfile content:
```
FROM golang:alpine
WORKDIR /app
COPY . .
CMD go run pingpong_demo/pong/pong.go
```

### Docker Image deployment
After creating the Dockerfiles, we must build the Docker Images and push each of them into the cluster, so they are available to be used in the Inspr environment.  

First, from "/pingpong_demo", go into "/ping" folder:
```
cd ping
```

Then we must build the Docker Image by using the Dockerfile previously created, applying the tag `latest` to it:
```
docker build -f Dockerfile -t YOUR-CLUSTER-URL/app/ping:latest
```

Finally, we push the builded image into the cluster:
```
docker push YOUR-CLUSTER-URL/app/ping:latest
```

Now the same must be done to Pong. Go back to "/pingpong_demo" folder, and access "/pong" folder:
```
cd ..
cd pong
```

Build the Docker Image by using the Dockerfile previously created, applying the tag `latest` to it:
```
docker build -f Dockerfile -t YOUR-CLUSTER-URL/app/pong:latest
```

Push the builded image into the cluster:
```
docker push YOUR-CLUSTER-URL/app/pong:latest
```  

**Alternatively, you can use Skaffold to do the steps described above, in a much simpler way.**  
If you have Skaffold installed, just go into the main folder "/pingpong_demo" and run the build command:
```
skaffold build
```
And it's done!