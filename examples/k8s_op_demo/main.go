package main

import (
	"context"
	"log"

	"inspr.dev/inspr/cmd/insprd/memory/tree"
	"inspr.dev/inspr/cmd/insprd/operators/kafka/nodes"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/utils"
)

func main() {
	mem := tree.GetTreeMemory()
	mem.InitTransaction()
	err := mem.ChannelTypes().Create("", &meta.ChannelType{
		Meta: meta.Metadata{
			Name: "channelType1",
		},
		Schema: "{\"type\":\"string\"}",
	})
	if err != nil {
		panic(err)
	}

	err = mem.Channels().Create("", &meta.Channel{
		Meta: meta.Metadata{
			Name: "ch1",
		},
		Spec: meta.ChannelSpec{
			Type: "channelType1",
		},
	})
	if err != nil {
		panic(err)
	}

	err = mem.Apps().Create("", &meta.App{
		Meta: meta.Metadata{
			Name:      "app1",
			Reference: "reference",
			Annotations: map[string]string{
				"app": "hellow",
			},
		},
		Spec: meta.AppSpec{
			Node: meta.Node{
				Spec: meta.NodeSpec{
					Image:    "gcr.io/red-inspr/inspr/sidecar/test:latest",
					Replicas: 4,
					Environment: utils.EnvironmentMap{
						"THIS_IS_AN_ENV_VAR": "THIS IS ITS VALUE",
					},
				},
			},
			Boundary: meta.AppBoundary{
				Input: utils.StringArray{
					"ch1",
				},
				Output: utils.StringArray{
					"ch1",
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	op, err := nodes.NewOperator(tree.GetTreeMemory())
	if err != nil {
		panic(err)
	}
	app, _ := mem.Apps().Get("app1")
	_, err = op.CreateNode(
		context.Background(),
		app,
	)
	if err != nil {
		log.Fatalf("%#v", err.Error())
	}

}
