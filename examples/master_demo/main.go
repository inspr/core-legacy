package main

import (
	"context"
	"fmt"
	"os"

	"inspr.dev/inspr/pkg/controller/client"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/rest/request"
)

func main() {
	client := client.Client{
		HTTPClient: request.NewJSONClient("http://127.0.0.1:8080"),
	}

	fmt.Println("[Creating App HelloWorld in Root...]")
	createHelloWorldApp(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Trying to create HelloWorld again...]")
	createHelloWorldApp(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Creating TypeHello inside HelloWorld app...]")
	createTypeInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Getting TypeHello...]")
	insprType, err := client.Types().Get(context.Background(), "HelloWorld", "TypeHello")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintTypeTree(insprType, os.Stdout)

	fmt.Println("[Creating ChannelOne inside HelloWorld app...]")
	createChannelInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Getting ChannelOne...]")
	ch, err := client.Channels().Get(context.Background(), "HelloWorld", "ChannelOne")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintChannelTree(ch, os.Stdout)

	fmt.Println("[Getting TypeHello...]")
	insprType, err = client.Types().Get(context.Background(), "HelloWorld", "TypeHello")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintTypeTree(insprType, os.Stdout)

	fmt.Println("[Creating NewApp inside HelloWorld app...]")
	createNewAppInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Update NewApp adding a new boundary: ChannelOne as Input...]")
	updateNewAppAddBoundary(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Update ChannelOne adding a annotaion to it...]")
	updateChannelOneAddAnnotationToIt(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Update TypeHello adding a note to it...]")
	updateTypeHelloAddAnnotation(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Delete NewApp inside HelloWorld...]")
	deleteNewAppInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Delete ChannelOne inside HelloWorld...]")
	deleteChannelOneInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Delete TypeHello inside HelloWorld]")
	deleteTypeHelloInsideHelloWorld(&client, false)
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
	utils.PrintAppTree(resp, os.Stdout)

	fmt.Println("[Creating NewApp inside HelloWorld app NOT USING DRY RUN...]")
	createNewAppInsideHelloWorld(&client, false)
	fmt.Printf("\n\n")

	fmt.Println("[Getting App HelloWorld...]")
	resp, err = client.Apps().Get(context.Background(), "HelloWorld")
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintAppTree(resp, os.Stdout)

	fmt.Println("[Deleting HelloWorld App...]")
	deleteHelloWorldApp(&client, false)
	fmt.Printf("\n\n")

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
	resp.Print(os.Stdout)
}

func createChannelInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Channels().Create(context.Background(), "HelloWorld", &meta.Channel{
		Meta: meta.Metadata{
			Name: "ChannelOne",
		},
		Spec: meta.ChannelSpec{
			Type: "TypeHello",
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print(os.Stdout)
}

func createTypeInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Types().Create(context.Background(), "HelloWorld", &meta.Type{
		Meta: meta.Metadata{
			Name: "TypeHello",
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print(os.Stdout)
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
	resp.Print(os.Stdout)
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
	resp.Print(os.Stdout)
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
			Type: "TypeHello",
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print(os.Stdout)
}

func updateTypeHelloAddAnnotation(client *client.Client, dryRun bool) {
	resp, err := client.Types().Update(context.Background(), "HelloWorld", &meta.Type{
		Meta: meta.Metadata{
			Name: "TypeHello",
			Annotations: map[string]string{
				"What's this?": "This is a note inside TypeHello",
			},
		},
	}, dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print(os.Stdout)
}

func deleteNewAppInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Apps().Delete(context.Background(), "HelloWorld.NewApp", dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print(os.Stdout)
}

func deleteChannelOneInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Channels().Delete(context.Background(), "HelloWorld", "ChannelOne", dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print(os.Stdout)
}

func deleteTypeHelloInsideHelloWorld(client *client.Client, dryRun bool) {
	resp, err := client.Types().Delete(context.Background(), "HelloWorld", "TypeHello", dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print(os.Stdout)
}

func deleteHelloWorldApp(client *client.Client, dryRun bool) {
	resp, err := client.Apps().Delete(context.Background(), "HelloWorld", dryRun)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	resp.Print(os.Stdout)
}
