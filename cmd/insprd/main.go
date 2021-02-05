// THIS IS THE MASTER
package main

import (
	"fmt"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/diff"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type person struct {
	name string
	age  int
}

func GetMockRootApp() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:        "",
			Reference:   "",
			Annotations: map[string]string{},
			Parent:      "",
			SHA256:      "1",
		},
		Spec: meta.AppSpec{
			Node: meta.Node{},
			Apps: map[string]*meta.App{
				"app1": {
					Meta: meta.Metadata{
						Name:        "app1",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp1",
								Reference:   "app1.nodeApp1",
								Annotations: map[string]string{},
								Parent:      "app1",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeApp1",
							},
						},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
				"app2": {
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "app2",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"app3": {
								Meta: meta.Metadata{
									Name:        "app3",
									Reference:   "app2.app3",
									Annotations: map[string]string{},
									Parent:      "app2",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "nodeApp3",
											Reference:   "app3.nodeApp2",
											Annotations: map[string]string{},
											Parent:      "app3",
											SHA256:      "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeApp3",
										},
									},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1"},
										Output: []string{"ch2"},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"ch1app2": {
								Meta: meta.Metadata{
									Name:   "ch1app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
							"ch2app2": {
								Meta: meta.Metadata{
									Name:   "ch2app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
			},
			Channels: map[string]*meta.Channel{
				"ch1": {},
				"ch2": {},
			},
			ChannelTypes: map[string]*meta.ChannelType{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "root.ct1",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
					Schema: []byte{},
				},
				"ct2": {
					Meta: meta.Metadata{
						Name:        "ct2",
						Reference:   "root.ct2",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
					Schema: []byte{},
				},
			},
			Boundary: meta.AppBoundary{
				Input:  []string{},
				Output: []string{},
			},
		},
	}
	return &root
}
func GetMockRootApp2() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:      "",
			Reference: "",
			Annotations: map[string]string{
				"an1": "alalalala",
				"an2": "lalalalal",
			},
			Parent: "",
			SHA256: "2",
		},
		Spec: meta.AppSpec{
			Node: meta.Node{},
			Apps: map[string]*meta.App{
				// "app1": {
				// 	Meta: meta.Metadata{
				// 		Name:        "app1",
				// 		Reference:   "blaa",
				// 		Annotations: map[string]string{},
				// 		Parent:      "",
				// 		SHA256:      "",
				// 	},
				// 	Spec: meta.AppSpec{
				// 		Node: meta.Node{
				// 			Meta: meta.Metadata{
				// 				Name:        "nodeApp1",
				// 				Reference:   "app1.nodeApp1",
				// 				Annotations: map[string]string{},
				// 				Parent:      "app1",
				// 				SHA256:      "",
				// 			},
				// 			Spec: meta.NodeSpec{
				// 				Image: "imageNodeApp1",
				// 			},
				// 		},
				// 		Apps:         map[string]*meta.App{},
				// 		Channels:     map[string]*meta.Channel{},
				// 		ChannelTypes: map[string]*meta.ChannelType{},
				// 		Boundary: meta.AppBoundary{
				// 			Input:  []string{"ch1"},
				// 			Output: []string{"ch2"},
				// 		},
				// 	},
				// },
				"app2": {
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "app2diff",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"app3": {
								Meta: meta.Metadata{
									Name:        "app3",
									Reference:   "app2.app3",
									Annotations: map[string]string{},
									Parent:      "app2",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "nodeApp3",
											Reference:   "app3.nodeApp2",
											Annotations: map[string]string{},
											Parent:      "app3",
											SHA256:      "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeApp3",
										},
									},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1"},
										Output: []string{"ch2"},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"ch1app2": {
								Meta: meta.Metadata{
									Name:   "ch1app2",
									Parent: "",
									SHA256: "diff",
								},
								Spec: meta.ChannelSpec{},
							},
							"ch2app2": {
								Meta: meta.Metadata{
									Name:   "ch2app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
			},
			Channels: map[string]*meta.Channel{
				"ch1": {},
				"ch2": {},
			},
			ChannelTypes: map[string]*meta.ChannelType{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "root.ct1",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
					Schema: []byte{},
				},
				"ct2": {
					Meta: meta.Metadata{
						Name:        "ct2",
						Reference:   "root.ct2",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
					Schema: []byte{},
				},
			},
			Boundary: meta.AppBoundary{
				Input:  []string{},
				Output: []string{},
			},
		},
	}
	return &root
}

func main() {
	// api.Run()
	tree1 := GetMockRootApp()
	tree2 := GetMockRootApp2()
	changelog, err := diff.Diff(tree1, tree2)
	if err != nil {
		fmt.Println(err.Error())
	}
	changelog.Print()
	// for _, change := range changelog {
	// 	for _, diff := range change.Diff {
	// 		fmt.Println("On: ", change.Context+diff.Field)
	// 		fmt.Println("From: ", diff.From)
	// 		fmt.Println("To: ", diff.Curr)
	// 		fmt.Println("")
	// 	}
	// }
}
