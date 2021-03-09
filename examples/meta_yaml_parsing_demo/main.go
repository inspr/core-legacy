package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	"gopkg.in/yaml.v2"
)

// tests if the app written in the /tmp/inspr-meta-complete-yaml-test.yaml
// when unmarshalled is equal to it's original definition
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
		log.Fatalln(err)
	}

	file, err := os.Create("/tmp/inspr-meta-complete-yaml-test.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	_, err = file.Write(bytesYAML)
	if err != nil {
		log.Fatalln(err)
	}

	var app meta.App

	data, _ := ioutil.ReadFile("/tmp/inspr-meta-complete-yaml-test.yaml")

	yaml.Unmarshal(data, &app)
	if changelog, _ := diff.Diff(appTest, &app); len(changelog) != 0 {
		log.Fatalln("TestServer_Init() = ", app, ", want ", appTest)
	}

	fmt.Println("worked as expected")
}
