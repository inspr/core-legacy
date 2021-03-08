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
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/fake"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func getMockApp() *meta.App {
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

func TestNewDescribeCmd(t *testing.T) {
	tests := []struct {
		name          string
		checkFunctiom func(t *testing.T, got *cobra.Command)
	}{
		{
			name: "It should create a new describe command",
			checkFunctiom: func(t *testing.T, got *cobra.Command) {
				if got == nil {
					t.Errorf("NewDescribeCmd() not created successfully")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDescribeCmd()
			if tt.checkFunctiom != nil {
				tt.checkFunctiom(t, got)
			}
		})
	}
}

func Test_getScope(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "It should return the scope",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cliutils.GetScope()
			if (err != nil) != tt.wantErr {
				t.Errorf("getScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getScope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processArg(t *testing.T) {
	type args struct {
		arg   string
		scope string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "Arg is a invalid structure name - it should return a error",
			args: args{
				arg:   "invalid!name",
				scope: "",
			},
			wantErr: true,
		},
		{
			name: "Arg is a valid structure name",
			args: args{
				arg:   "helloWorld",
				scope: "app1",
			},
			want:    "app1",
			want1:   "helloWorld",
			wantErr: false,
		},
		{
			name: "Arg is a invalid scope - it should return a error",
			args: args{
				arg:   "hello..World",
				scope: "app1",
			},
			wantErr: true,
		},
		{
			name: "Arg is a valid scope",
			args: args{
				arg:   "hello.World",
				scope: "app1",
			},
			want:    "app1.hello",
			want1:   "World",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := cliutils.ProcessArg(tt.args.arg, tt.args.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("processArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("processArg() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("processArg() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_displayAppState(t *testing.T) {
	bufResp := bytes.NewBufferString("")
	utils.PrintAppTree(getMockApp(), bufResp)
	outResp, _ := ioutil.ReadAll(bufResp)

	ah := fake.MockMemoryManager(nil)
	ah.Apps().CreateApp("", getMockApp())

	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		app, err := ah.Apps().GetApp(data.Ctx)
		if err != nil {
			fmt.Println(err)
		}
		rest.JSON(w, http.StatusOK, app)
	}

	tests := []struct {
		name           string
		flagsAndArgs   []string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
		expectedOutput []byte
	}{
		{
			name:           "Should describe the app state",
			flagsAndArgs:   []string{"a", "appParent"},
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
			flagsAndArgs:   []string{"a", "", "--scope", "appParent"},
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
			cmd := NewDescribeCmd()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			cliutils.SetClient(server.URL)

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("displayAppState() = %v, want %v", string(got), string(tt.expectedOutput))
			}
		})
	}
}

func Test_displayChannelState(t *testing.T) {
	bufResp := bytes.NewBufferString("")
	utils.PrintChannelTree(getMockApp().Spec.Channels["ch1"], bufResp)
	outResp, _ := ioutil.ReadAll(bufResp)

	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		ch := getMockApp().Spec.Channels[data.ChName]

		rest.JSON(w, http.StatusOK, ch)
	}

	tests := []struct {
		name           string
		flagsAndArgs   []string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
		expectedOutput []byte
	}{
		{
			name:           "Should describe the channel state",
			flagsAndArgs:   []string{"ch", "appParent.ch1"},
			handlerFunc:    handler,
			expectedOutput: outResp,
		},
		{
			name:           "Invalid scope flag, should not print",
			flagsAndArgs:   []string{"ch", "ch1", "--scope", "invalid..scope"},
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
			flagsAndArgs:   []string{"ch", "invalid..args", "--scope", "appParent"},
			handlerFunc:    handler,
			expectedOutput: []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDescribeCmd()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			cliutils.SetClient(server.URL)

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("displayChannelState() = %v, want %v", string(got), string(tt.expectedOutput))
			}
		})
	}
}

func Test_displayChannelTypeState(t *testing.T) {
	bufResp := bytes.NewBufferString("")
	utils.PrintChannelTypeTree(getMockApp().Spec.ChannelTypes["ct1"], bufResp)
	outResp, _ := ioutil.ReadAll(bufResp)

	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		ct := getMockApp().Spec.ChannelTypes[data.CtName]

		rest.JSON(w, http.StatusOK, ct)
	}

	tests := []struct {
		name           string
		flagsAndArgs   []string
		handlerFunc    func(w http.ResponseWriter, r *http.Request)
		expectedOutput []byte
	}{
		{
			name:           "Should describe the channelType state",
			flagsAndArgs:   []string{"ct", "appParent.ct1"},
			handlerFunc:    handler,
			expectedOutput: outResp,
		},
		{
			name:           "Invalid scope flag, should not print",
			flagsAndArgs:   []string{"ct", "ct1", "--scope", "invalid..scope"},
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
			flagsAndArgs:   []string{"ct", "invalid..args", "--scope", "appParent"},
			handlerFunc:    handler,
			expectedOutput: []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDescribeCmd()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			server := httptest.NewServer(http.HandlerFunc(tt.handlerFunc))
			cliutils.SetClient(server.URL)

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf("displayChannelTypeState() = %v, want %v", string(got), string(tt.expectedOutput))
			}
		})
	}
}
