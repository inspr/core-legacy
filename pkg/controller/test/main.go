package main

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

func main() {
	client := client.Client{
		HTTPClient: request.NewClient().BaseURL("http://127.0.0.1:8080").Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build(),
	}

	fmt.Println("[Creating App HelloWorld in Root...]")
	createHelloWorldApp(&client)
	fmt.Printf("\n\n")

	fmt.Println("[Trying to create HelloWorld again...]")
	createHelloWorldApp(&client)
	fmt.Printf("\n\n")

	fmt.Println("[Creating ChannelTypeHello inside HelloWorld app...]")
	createChannelTypeInsideHelloWorld(&client)
	fmt.Printf("\n\n")

	fmt.Println("[Creating ChannelOne inside HelloWorld app...]")
	createChannelInsideHelloWorld(&client)
	fmt.Printf("\n\n")

	fmt.Println("[Creating NewApp inside HelloWorld app...]")
	createNewAppInsideHelloWorld(&client)
	fmt.Printf("\n\n")

	fmt.Println("[Update NewApp adding a new boundary: ChannelOne as Input...]")
	updateNewAppAddBoundary(&client)
	fmt.Printf("\n\n")

}

func createHelloWorldApp(client *client.Client) {
	resp, err := client.Apps().Create(context.Background(), "", &meta.App{
		Meta: meta.Metadata{
			Name: "HelloWorld",
		},
		Spec: meta.AppSpec{},
	})
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func createChannelInsideHelloWorld(client *client.Client) {
	resp, err := client.Channels().Create(context.Background(), "HelloWorld", &meta.Channel{
		Meta: meta.Metadata{
			Name: "ChannelOne",
		},
		Spec: meta.ChannelSpec{
			Type: "ChannelTypeHello",
		},
	})
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func createChannelTypeInsideHelloWorld(client *client.Client) {
	resp, err := client.ChannelTypes().Create(context.Background(), "HelloWorld", &meta.ChannelType{
		Meta: meta.Metadata{
			Name: "ChannelTypeHello",
		},
	})
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func createNewAppInsideHelloWorld(client *client.Client) {
	resp, err := client.Apps().Create(context.Background(), "HelloWorld", &meta.App{
		Meta: meta.Metadata{
			Name: "NewApp",
		},
		Spec: meta.AppSpec{},
	})
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func updateNewAppAddBoundary(client *client.Client) {
	resp, err := client.Apps().Update(context.Background(), "HelloWorld.NewApp", &meta.App{
		Meta: meta.Metadata{
			Name: "NewApp",
		},
		Spec: meta.AppSpec{
			Boundary: meta.AppBoundary{
				Input: []string{"ChannelOne"},
			},
		},
	})
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}
