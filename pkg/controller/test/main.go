package main

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
	"gitlab.inspr.dev/inspr/core/pkg/utils/diff"
)

func main() {
	client := client.Client{
		HTTPClient: request.NewClient().BaseURL("http://127.0.0.1:8080").Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build(),
	}

	fmt.Println("[Creating App HelloWorld in Root...]")
	resp, _ := createHelloWorldApp(&client)
	resp.Print()
	fmt.Printf("\n\n")

	fmt.Println("[Trying to create HelloWorld again...]")
	_, err := createHelloWorldApp(&client)
	fmt.Println(err)
	fmt.Printf("\n\n")

	fmt.Println("[Trying to create ChannelOne inside HelloWorld app...]")
	resp, _ = createChannelInsideHelloWorld(&client)
	resp.Print()
	fmt.Printf("\n\n")

}

func createHelloWorldApp(client *client.Client) (diff.Changelog, error) {
	resp, err := client.Apps().Create(context.Background(), "", &meta.App{
		Meta: meta.Metadata{
			Name: "HelloWorld",
		},
		Spec: meta.AppSpec{
			ChannelTypes: map[string]*meta.ChannelType{
				"ChannelTypeHello": {
					Meta: meta.Metadata{
						Name: "ChannelTypeHello",
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}
	return resp, nil
}

func createChannelInsideHelloWorld(client *client.Client) (diff.Changelog, error) {
	resp, err := client.Channels().Create(context.Background(), "HelloWorld", &meta.Channel{
		Meta: meta.Metadata{
			Name: "channelOne",
		},
		Spec: meta.ChannelSpec{
			Type: "ChannelTypeHello",
		},
	})
	if err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}
	return resp, nil
}
