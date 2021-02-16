package main

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/controller/client"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

func main() {
	client := client.Client{
		HTTPClient: request.NewClient().BaseURL("http://127.0.0.1:8080").Encoder(json.Marshal).Decoder(request.JSONDecoderGenerator).Build(),
	}

	fmt.Println("[Creating App HelloWorld in Root...]")
	createHelloWorldApp(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Trying to create HelloWorld again...]")
	createHelloWorldApp(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Creating ChannelTypeHello inside HelloWorld app...]")
	createChannelTypeInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Getting ChannelTypeHello...]")
	cType, err := client.ChannelTypes().Get(context.Background(), "HelloWorld", "ChannelTypeHello")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintChannelTypeTree(cType)

	fmt.Println("[Creating ChannelOne inside HelloWorld app...]")
	createChannelInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Getting ChannelOne...]")
	ch, err := client.Channels().Get(context.Background(), "HelloWorld", "ChannelOne")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintChannelTree(ch)

	fmt.Println("[Getting ChannelTypeHello...]")
	cType, err = client.ChannelTypes().Get(context.Background(), "HelloWorld", "ChannelTypeHello")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintChannelTypeTree(cType)

	fmt.Println("[Creating NewApp inside HelloWorld app...]")
	createNewAppInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Update NewApp adding a new boundary: ChannelOne as Input...]")
	updateNewAppAddBoundary(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Update ChannelOne adding a annotaion to it...]")
	updateChannelOneAddAnnotationToIt(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Update ChannelTypeHello adding a note to it...]")
	updateChannelTypeHelloAddAnnotation(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Delete NewApp inside HelloWorld...]")
	deleteNewAppInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Delete ChannelOne inside HelloWorld...]")
	deleteChannelOneInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Delete ChannelTypeHello inside HelloWorld]")
	deleteChannelTypeHelloInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[TESTING DRY RUN...]")
	fmt.Printf("\n\n")

	fmt.Println("[Creating NewApp inside HelloWorld app USING DRY RUN...]")
	createNewAppInsideHelloWorld(&client, true)
	fmt.Printf("\n\n")

	fmt.Println("[Getting App HelloWorld...]")
	resp, err := client.Apps().Get(context.Background(), "HelloWorld")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintAppTree(resp)

	fmt.Println("[Creating NewApp inside HelloWorld app NOT USING DRY RUN...]")
	createNewAppInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Getting App HelloWorld...]")
	resp, err = client.Apps().Get(context.Background(), "HelloWorld")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintAppTree(resp)

}

func createHelloWorldApp(client *client.Client, dryRun bool) {
	resp, err := client.Apps().Create(context.Background(), "", &meta.App{
		Meta: meta.Metadata{
			Name: "HelloWorld",
		},
		Spec: meta.AppSpec{},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func createChannelInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Channels().Create(context.Background(), "HelloWorld", &meta.Channel{
		Meta: meta.Metadata{
			Name: "ChannelOne",
		},
		Spec: meta.ChannelSpec{
			Type: "ChannelTypeHello",
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func createChannelTypeInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.ChannelTypes().Create(context.Background(), "HelloWorld", &meta.ChannelType{
		Meta: meta.Metadata{
			Name: "ChannelTypeHello",
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func createNewAppInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Apps().Create(context.Background(), "HelloWorld", &meta.App{
		Meta: meta.Metadata{
			Name: "NewApp",
		},
		Spec: meta.AppSpec{},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func updateNewAppAddBoundary(client *client.Client, dryRun bool) {
	resp, err := client.Apps().Update(context.Background(), "HelloWorld.NewApp", &meta.App{
		Meta: meta.Metadata{
			Name: "NewApp",
		},
		Spec: meta.AppSpec{
			Boundary: meta.AppBoundary{
				Input: []string{"ChannelOne"},
			},
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func updateChannelOneAddAnnotationToIt(client *client.Client, dryRun bool) {
	resp, err := client.Channels().Update(context.Background(), "HelloWorld", &meta.Channel{
		Meta: meta.Metadata{
			Name: "ChannelOne",
			Annotations: map[string]string{
				"NoteOne": "A brand new note!",
			},
		},
		Spec: meta.ChannelSpec{
			Type: "ChannelTypeHello",
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func updateChannelTypeHelloAddAnnotation(client *client.Client, dryRun bool) {
	resp, err := client.ChannelTypes().Update(context.Background(), "HelloWorld", &meta.ChannelType{
		Meta: meta.Metadata{
			Name: "ChannelTypeHello",
			Annotations: map[string]string{
				"What's this?": "This is a note inside ChannelTypeHello",
			},
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func deleteNewAppInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Apps().Delete(context.Background(), "HelloWorld.NewApp", dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func deleteChannelOneInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Channels().Delete(context.Background(), "HelloWorld", "ChannelOne", dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}

func deleteChannelTypeHelloInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.ChannelTypes().Delete(context.Background(), "HelloWorld", "ChannelTypeHello", dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print()
}
