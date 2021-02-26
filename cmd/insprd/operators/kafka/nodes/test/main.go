package main

import (
	"context"
	"log"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators/kafka/nodes"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

func main() {
	mem := tree.GetTreeMemory()
	mem.InitTransaction()
	err := mem.ChannelTypes().CreateChannelType("", &meta.ChannelType{
		Meta: meta.Metadata{
			Name: "channelType1",
		},
		Schema: "{\"type\":\"string\"}",
	})
	if err != nil {
		panic(err)
	}

	err = mem.Channels().CreateChannel("", &meta.Channel{
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

	err = mem.Apps().CreateApp("", &meta.App{
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
					Replicas: 3,
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

	op, err := nodes.NewOperator()
	if err != nil {
		panic(err)
	}
	app, _ := mem.Apps().Get("app1")
	_, err = op.CreateNode(
		context.Background(),
		app,
	)
	if err != nil {
		log.Fatalf("%#v", err.(*ierrors.InsprError).Err.Error())
	}

}
