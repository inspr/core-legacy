package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func getMockAppWithoutApp1() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:        "appParent",
			Reference:   "appParent",
			Annotations: map[string]string{},
			Parent:      "",
			SHA256:      "",
		},
		Spec: meta.AppSpec{
			Node: meta.Node{},
			Apps: map[string]*meta.App{},
			Channels: map[string]*meta.Channel{
				"ch1": {
					Meta: meta.Metadata{
						Name:   "ch1",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "ct1",
					},
				},
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
					Schema: "",
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

func getMockAppWithoutCh1() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:        "appParent",
			Reference:   "appParent",
			Annotations: map[string]string{},
			Parent:      "",
			SHA256:      "",
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
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"thenewapp": {
								Meta: meta.Metadata{
									Name:        "thenewapp",
									Reference:   "app1.thenewapp",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1app1"},
										Output: []string{},
									},
								},
							},
						},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch1"},
						},
					},
				},
			},
			Channels: map[string]*meta.Channel{
				"ch1": {
					Meta: meta.Metadata{
						Name:   "ch1",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "ct1",
					},
				},
			},
			ChannelTypes: map[string]*meta.ChannelType{},
			Boundary: meta.AppBoundary{
				Input:  []string{},
				Output: []string{},
			},
		},
	}
	return &root
}

func getMockAppWithoutCt1() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:        "appParent",
			Reference:   "appParent",
			Annotations: map[string]string{},
			Parent:      "",
			SHA256:      "",
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
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"thenewapp": {
								Meta: meta.Metadata{
									Name:        "thenewapp",
									Reference:   "app1.thenewapp",
									Annotations: map[string]string{},
									Parent:      "app1",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1app1"},
										Output: []string{},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"ch1app1": {
								Meta: meta.Metadata{
									Name:   "ch1app1",
									Parent: "",
								},
								ConnectedApps: []string{"thenewapp"},
								Spec:          meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch1"},
						},
					},
				},
			},
			Channels: map[string]*meta.Channel{
				"ch1": {
					Meta: meta.Metadata{
						Name:   "ch1",
						Parent: "",
					},
					Spec: meta.ChannelSpec{
						Type: "ct1",
					},
				},
			},
			ChannelTypes: map[string]*meta.ChannelType{},
			Boundary: meta.AppBoundary{
				Input:  []string{},
				Output: []string{},
			},
		},
	}
	return &root
}

func TestNewDeleteCmd(t *testing.T) {
	defer restartScopeFlag()
	tests := []struct {
		name          string
		checkFunction func(t *testing.T, got *cobra.Command)
	}{
		{
			name: "It should create a new delete command",
			checkFunction: func(t *testing.T, got *cobra.Command) {
				if got == nil {
					t.Errorf("NewDeleteCmd() not created successfully")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDeleteCmd()
			if tt.checkFunction != nil {
				tt.checkFunction(t, got)
			}
		})
	}
}

func Test_deleteApps(t *testing.T) {
	defer restartScopeFlag()
	bufResp := bytes.NewBufferString("")
	changelog, _ := diff.Diff(getMockApp(), getMockAppWithoutApp1())

	changelog.Print(bufResp)
	outResp, _ := ioutil.ReadAll(bufResp)

	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		if data.Ctx != "appParent.app1" {
			rest.ERROR(w, ierrors.NewError().Message("error test").Build())
			return
		}

		rest.JSON(w, http.StatusOK, changelog)
	}

	tests := []struct {
		name           string
		flagsAndArgs   []string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
		expectedOutput []byte
	}{
		{
			name:           "Should delete the app and return the diff",
			flagsAndArgs:   []string{"a", "appParent.app1"},
			handlerFunc:    handler,
			expectedOutput: outResp,
		},
		{
			name:           "Invalid scope flag, should not print",
			flagsAndArgs:   []string{"a", "appParent", "--scope", "invalid..scope"},
			handlerFunc:    handler,
			expectedOutput: []byte(""),
		},
		{
			name:           "Valid scope flag",
			flagsAndArgs:   []string{"a", "", "--scope", "appParent.app1"},
			handlerFunc:    handler,
			expectedOutput: outResp,
		},
		{
			name:           "Invalid arg",
			flagsAndArgs:   []string{"a", "invalid..args", "--scope", "appParent"},
			handlerFunc:    handler,
			expectedOutput: []byte("invalid args\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDeleteCmd()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			cliutils.SetClient(server.URL)

			defer server.Close()

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("deleteApps() = %v, want %v", string(got), string(tt.expectedOutput))
			}
		})
	}
}

func Test_deleteChannels(t *testing.T) {
	defer restartScopeFlag()
	bufResp := bytes.NewBufferString("")
	changelog, _ := diff.Diff(getMockApp(), getMockAppWithoutCh1())

	changelog.Print(bufResp)
	outResp, _ := ioutil.ReadAll(bufResp)

	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		if data.Ctx != "appParent" || data.ChName != "ch1" {
			rest.ERROR(w, ierrors.NewError().Message("error test").Build())
			return
		}

		rest.JSON(w, http.StatusOK, changelog)
	}

	tests := []struct {
		name           string
		flagsAndArgs   []string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
		expectedOutput []byte
	}{
		{
			name:           "Should delete the channel and return the diff",
			flagsAndArgs:   []string{"ch", "appParent.ch1"},
			handlerFunc:    handler,
			expectedOutput: outResp,
		},
		{
			name:           "Invalid scope flag, should not print",
			flagsAndArgs:   []string{"ch", "appParent.ch1", "--scope", "invalid..scope"},
			handlerFunc:    handler,
			expectedOutput: []byte(""),
		},
		{
			name:           "Valid scope flag",
			flagsAndArgs:   []string{"ch", "ch1", "--scope", "appParent"},
			handlerFunc:    handler,
			expectedOutput: outResp,
		},
		{
			name:           "Invalid arg",
			flagsAndArgs:   []string{"ch", "invalid..args", "--scope", "appParent.ch1"},
			handlerFunc:    handler,
			expectedOutput: []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDeleteCmd()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			cliutils.SetClient(server.URL)

			defer server.Close()

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("deleteChannels() = %v, want %v", string(got), string(tt.expectedOutput))
			}
		})
	}
}

func Test_deleteCTypes(t *testing.T) {
	defer restartScopeFlag()
	bufResp := bytes.NewBufferString("")
	changelog, _ := diff.Diff(getMockApp(), getMockAppWithoutCt1())

	changelog.Print(bufResp)
	outResp, _ := ioutil.ReadAll(bufResp)

	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		if data.Ctx != "appParent" || data.CtName != "ct1" {
			rest.ERROR(w, ierrors.NewError().Message("error test").Build())
			return
		}

		rest.JSON(w, http.StatusOK, changelog)
	}

	tests := []struct {
		name           string
		flagsAndArgs   []string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
		expectedOutput []byte
	}{
		{
			name:           "Should delete the channelType and return the diff",
			flagsAndArgs:   []string{"ct", "appParent.ct1"},
			handlerFunc:    handler,
			expectedOutput: outResp,
		},
		{
			name:           "Invalid scope flag, should not print",
			flagsAndArgs:   []string{"ct", "appParent.ct1", "--scope", "invalid..scope"},
			handlerFunc:    handler,
			expectedOutput: []byte(""),
		},
		{
			name:           "Valid scope flag",
			flagsAndArgs:   []string{"ct", "ct1", "--scope", "appParent"},
			handlerFunc:    handler,
			expectedOutput: outResp,
		},
		{
			name:           "Invalid arg",
			flagsAndArgs:   []string{"ct", "invalid..args", "--scope", "appParent.ct1"},
			handlerFunc:    handler,
			expectedOutput: []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDeleteCmd()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			cliutils.SetClient(server.URL)

			defer server.Close()

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("deleteCtypes() = %v, want %v", string(got), string(tt.expectedOutput))
			}
		})
	}
}