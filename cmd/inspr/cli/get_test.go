package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"text/tabwriter"

	"github.com/spf13/cobra"
	cliutils "gitlab.inspr.dev/inspr/core/cmd/inspr/cli/utils"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

func TestNewGetCmd(t *testing.T) {
	tests := []struct {
		name          string
		checkFunction func(t *testing.T, got *cobra.Command)
	}{
		{
			name: "It should create a new describe command",
			checkFunction: func(t *testing.T, got *cobra.Command) {
				if got == nil {
					t.Errorf("NewDescribeCmd() not created successfully")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDescribeCmd()
			if tt.checkFunction != nil {
				tt.checkFunction(t, got)
			}
		})
	}
}

func Test_getApps(t *testing.T) {
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "appParent\n")
	fmt.Fprint(tabWriter, "app1\n")
	fmt.Fprint(tabWriter, "thenewapp\n")
	tabWriter.Flush()
	outResp, _ := ioutil.ReadAll(bufResp)
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		tab     []byte
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getApps valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     outResp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				data := models.AppQueryDI{}
				decoder := json.NewDecoder(r.Body)

				err := decoder.Decode(&data)
				if err != nil {
					fmt.Println(err)
				}

				app := getMockApp()

				rest.JSON(w, http.StatusOK, app)
			},
		},
		{
			name: "getApps invalid test, HTTP error",
			args: args{
				in0: context.Background(),
			},
			wantErr: true,
			tab:     outResp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.NewError().Message("error").Build())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getApps(tt.args.in0)
			got, _ := ioutil.ReadAll(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("getApps() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("getApps() error = %v, want %v", string(got), string(tt.tab))
			}
		})
	}
}

func Test_getChannels(t *testing.T) {
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "ch1\n")
	fmt.Fprint(tabWriter, "ch1app1\n")
	tabWriter.Flush()
	outResp, _ := ioutil.ReadAll(bufResp)
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		tab     []byte
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getChannels valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     outResp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				data := models.AppQueryDI{}
				decoder := json.NewDecoder(r.Body)

				err := decoder.Decode(&data)
				if err != nil {
					fmt.Println(err)
				}

				app := getMockApp()

				rest.JSON(w, http.StatusOK, app)
			},
		},
		{
			name: "getChannels invalid test, HTTP error",
			args: args{
				in0: context.Background(),
			},
			wantErr: true,
			tab:     outResp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.NewError().Message("error").Build())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getChannels(tt.args.in0)
			got, _ := ioutil.ReadAll(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("getChannels() error = %v, want %v", string(got), string(tt.tab))
			}
		})
	}
}

func Test_getCTypes(t *testing.T) {
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "ct1\n")
	tabWriter.Flush()
	outResp, _ := ioutil.ReadAll(bufResp)
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		tab     []byte
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getCTypes valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     outResp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				data := models.AppQueryDI{}
				decoder := json.NewDecoder(r.Body)

				err := decoder.Decode(&data)
				if err != nil {
					fmt.Println(err)
				}

				app := getMockApp()

				rest.JSON(w, http.StatusOK, app)
			},
		},
		{
			name: "getCTypes invalid test, HTTP error",
			args: args{
				in0: context.Background(),
			},
			wantErr: true,
			tab:     outResp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.NewError().Message("error").Build())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getCTypes(tt.args.in0)
			got, _ := ioutil.ReadAll(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("getCTypes() error = %v, want %v", string(got), string(tt.tab))
			}
		})
	}
}

func Test_getNodes(t *testing.T) {

	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "thenewapp\n")
	tabWriter.Flush()
	outResp, _ := ioutil.ReadAll(bufResp)
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		tab     []byte
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getNodes valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     outResp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				data := models.AppQueryDI{}
				decoder := json.NewDecoder(r.Body)

				err := decoder.Decode(&data)
				if err != nil {
					fmt.Println(err)
				}

				app := getMockApp()

				rest.JSON(w, http.StatusOK, app)
			},
		},
		{
			name: "getNodes invalid test, HTTP error",
			args: args{
				in0: context.Background(),
			},
			wantErr: true,
			tab:     outResp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.NewError().Message("error").Build())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getNodes(tt.args.in0)
			got, _ := ioutil.ReadAll(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNodes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("getNodes() error = %v, want %v", string(got), string(tt.tab))
			}
		})
	}
}

func Test_getObj(t *testing.T) {
	type args struct {
		printObj func(*meta.App, *[]string)
		lines    *[]string
	}
	tests := []struct {
		name    string
		args    args
		lines   *[]string
		wantErr bool
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getObj valid execution, printApps case",
			args: args{
				printObj: printApps,
				lines:    &[]string{},
			},
			lines:   &[]string{"appParent\n", "app1\n", "thenewapp\n"},
			wantErr: false,
			handler: func(w http.ResponseWriter, r *http.Request) {
				data := models.AppQueryDI{}
				decoder := json.NewDecoder(r.Body)

				err := decoder.Decode(&data)
				if err != nil {
					fmt.Println(err)
				}

				app := getMockApp()

				rest.JSON(w, http.StatusOK, app)
			},
		},
		{
			name: "getObj invalid execution, HTTP error",
			args: args{
				printObj: printApps,
				lines:    &[]string{},
			},
			lines:   &[]string{"appParent\n", "app1\n", "thenewapp\n"},
			wantErr: true,
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.NewError().Message("error").Build())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL)
			if err := getObj(tt.args.printObj, tt.args.lines); (err != nil) != tt.wantErr {
				t.Errorf("getObj() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("getObj() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_printApps(t *testing.T) {
	type args struct {
		app   *meta.App
		lines *[]string
	}
	tests := []struct {
		name  string
		args  args
		lines *[]string
	}{
		{
			name: "printApps test",
			args: args{
				app:   getMockApp(),
				lines: &[]string{},
			},
			lines: &[]string{"appParent\n", "app1\n", "thenewapp\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printApps(tt.args.app, tt.args.lines)
			if !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("printApps() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_printChannels(t *testing.T) {
	type args struct {
		app   *meta.App
		lines *[]string
	}
	tests := []struct {
		name  string
		args  args
		lines *[]string
	}{
		{
			name: "printChannels test",
			args: args{
				app:   getMockApp(),
				lines: &[]string{},
			},
			lines: &[]string{"ch1\n", "ch1app1\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printChannels(tt.args.app, tt.args.lines)
			if !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("printChannels() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_printCTypes(t *testing.T) {
	type args struct {
		app   *meta.App
		lines *[]string
	}
	tests := []struct {
		name  string
		args  args
		lines *[]string
	}{
		{
			name: "printCTypes test",
			args: args{
				app:   getMockApp(),
				lines: &[]string{},
			},
			lines: &[]string{"ct1\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printCTypes(tt.args.app, tt.args.lines)
			if !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("printCTypes() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_printNodes(t *testing.T) {
	type args struct {
		app   *meta.App
		lines *[]string
	}
	tests := []struct {
		name  string
		args  args
		lines *[]string
	}{
		{
			name: "printNodes test",
			args: args{
				app:   getMockApp(),
				lines: &[]string{},
			},
			lines: &[]string{"thenewapp\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printNodes(tt.args.app, tt.args.lines)
			if !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("printNodes() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_printLine(t *testing.T) {
	type args struct {
		name  string
		lines *[]string
	}
	tests := []struct {
		name  string
		args  args
		lines *[]string
	}{
		{
			name: "printLine test",
			args: args{
				name:  "line",
				lines: &[]string{},
			},
			lines: &[]string{"line\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printLine(tt.args.name, tt.args.lines)
			if !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("printLine() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_initTab(t *testing.T) {
	type args struct {
		lines *[]string
	}
	tests := []struct {
		name  string
		args  args
		lines *[]string
	}{
		{
			name: "initTab test",
			args: args{
				lines: &[]string{},
			},
			lines: &[]string{"NAME\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTab(tt.args.lines)
			if !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("initTab() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_printTab(t *testing.T) {
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)

	fmt.Fprint(tabWriter, "line\n")
	tabWriter.Flush()

	outResp, _ := ioutil.ReadAll(bufResp)

	type args struct {
		lines *[]string
	}
	tests := []struct {
		name string
		args args
		tab  []byte
	}{
		{
			name: "printtab test",
			args: args{
				lines: &[]string{"line\n"},
			},
			tab: outResp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			printTab(tt.args.lines)
			got, _ := ioutil.ReadAll(buf)

			if !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("printTab() error = %v, want %v", string(got), string(tt.tab))
			}
		})
	}
}