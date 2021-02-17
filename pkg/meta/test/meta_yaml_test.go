package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
	"gopkg.in/yaml.v2"
)

func main() {
	appTest := &meta.App{
		Meta: meta.Metadata{
			Name: "App1",
		},
		Spec: meta.AppSpec{
			Apps: map[string]*meta.App{
				"app2": {
					Meta: meta.Metadata{
						Name:      "app2",
						Reference: "App1.app2",
					},
					Spec: meta.AppSpec{
						Channels: map[string]*meta.Channel{
							"channel1": {
								Meta: meta.Metadata{
									Name: "channel1",
								},
								Spec: meta.ChannelSpec{
									Type: "ct1",
								},
							},
							"channel2": {
								Meta: meta.Metadata{
									Name: "channel2",
								},
								Spec: meta.ChannelSpec{
									Type: "ct1",
								},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
							"ct1": {
								Meta: meta.Metadata{
									Name: "ct1",
								},
								ConnectedChannels: []string{"channel1", "channel2"},
							},
						},
					},
				},
				"app3": {
					Meta: meta.Metadata{
						Name:      "app3",
						Reference: "App1.app3",
					},
					Spec: meta.AppSpec{
						Channels: map[string]*meta.Channel{
							"channel1": {
								Meta: meta.Metadata{
									Name: "channel1",
								},
								Spec: meta.ChannelSpec{
									Type: "ct1",
								},
							},
							"channel2": {
								Meta: meta.Metadata{
									Name: "channel2",
								},
								Spec: meta.ChannelSpec{
									Type: "ct1",
								},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{
							"ct1": {
								Meta: meta.Metadata{
									Name: "ct1",
								},
								ConnectedChannels: []string{"channel1", "channel2"},
							},
						},
					},
				},
			},
			Channels: map[string]*meta.Channel{
				"channel2": {
					Meta: meta.Metadata{
						Name: "channel2",
					},
				},
			},
			Boundary: meta.AppBoundary{
				Input:  []string{"ch1", "ch2"},
				Output: []string{"ch1", "ch2", "ch3"},
			},
		},
	}

	bytesYAML, err := yaml.Marshal(appTest)

	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create("/tmp/averyuniquename.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = file.Write(bytesYAML)
	if err != nil {
		fmt.Println(err)
		return
	}

	var app meta.App

	data, err := ioutil.ReadFile("/tmp/averyuniquename.yaml")

	yaml.Unmarshal(data, &app)

	utils.PrintAppTree(&app)

	app2 := app.Spec.Apps["app2"]
	utils.PrintAppTree(app2)

	channel2 := app2.Spec.Channels["channel2"]
	utils.PrintChannelTree(channel2)

	ct1 := app2.Spec.ChannelTypes["ct1"]
	utils.PrintChannelTypeTree(ct1)
}
