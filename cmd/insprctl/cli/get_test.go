package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"inspr.dev/inspr/pkg/api/models"
	cliutils "inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/rest"
)

func TestNewGetCmd(t *testing.T) {
	prepareToken(t)
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
	prepareToken(t)
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "appParent\n")
	fmt.Fprint(tabWriter, "app1\n")
	fmt.Fprint(tabWriter, "thenewapp\n")
	tabWriter.Flush()

	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		host    string
		args    args
		wantErr bool
		tab     string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getApps valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     bufResp.String(),
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
			tab:     bufResp.String(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL, tt.host)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getApps(tt.args.in0)
			got := buf.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("getApps() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("getApps() error = %v, want %v", got, tt.tab)
			}
		})
	}
}

func Test_getChannels(t *testing.T) {
	prepareToken(t)
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "ch1\n")
	fmt.Fprint(tabWriter, "ch1app1\n")
	tabWriter.Flush()

	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		host    string
		args    args
		wantErr bool
		tab     string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getChannels valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     bufResp.String(),
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
			tab:     bufResp.String(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL, tt.host)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getChannels(tt.args.in0)
			got := buf.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("getChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("getChannels() error = %v, want %v", got, tt.tab)
			}
		})
	}
}

func Test_gettypes(t *testing.T) {
	prepareToken(t)
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "ct1\n")
	tabWriter.Flush()

	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		host    string
		args    args
		wantErr bool
		tab     string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "gettypes valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     bufResp.String(),
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
			name: "gettypes invalid test, HTTP error",
			args: args{
				in0: context.Background(),
			},
			wantErr: true,
			tab:     bufResp.String(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL, tt.host)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getTypes(tt.args.in0)
			got := buf.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("gettypes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("gettypes() error = %v, want %v", got, tt.tab)
			}
		})
	}
}

func Test_getNodes(t *testing.T) {
	prepareToken(t)

	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "thenewapp\n")
	tabWriter.Flush()

	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		host    string
		args    args
		wantErr bool
		tab     string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getNodes valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     bufResp.String(),
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
			tab:     bufResp.String(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL, tt.host)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getNodes(tt.args.in0)
			got := buf.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("getNodes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("getNodes() error = %v, want %v", got, tt.tab)
			}
		})
	}
}

func Test_getObj(t *testing.T) {
	prepareToken(t)
	type args struct {
		printObj func(*meta.App, *[]string)
		lines    *[]string
	}
	tests := []struct {
		name    string
		host    string
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
				rest.ERROR(w, ierrors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL, tt.host)

			client := cliutils.GetCliClient()
			out := cliutils.GetCliOutput()
			scope, _ := cliutils.GetScope()

			if err := getObj(tt.args.printObj, tt.args.lines, client, out, scope); (err != nil) != tt.wantErr {
				t.Errorf("getObj() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("getObj() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_printApps(t *testing.T) {
	prepareToken(t)
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
	prepareToken(t)
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

func Test_printtypes(t *testing.T) {
	prepareToken(t)
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
			name: "printtypes test",
			args: args{
				app:   getMockApp(),
				lines: &[]string{},
			},
			lines: &[]string{"ct1\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printTypes(tt.args.app, tt.args.lines)
			if !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("printtypes() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_printNodes(t *testing.T) {
	prepareToken(t)
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
	prepareToken(t)
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
	prepareToken(t)
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
	prepareToken(t)
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)

	fmt.Fprint(tabWriter, "line\n")
	tabWriter.Flush()

	type args struct {
		lines *[]string
	}
	tests := []struct {
		name string
		args args
		tab  string
	}{
		{
			name: "printtab test",
			args: args{
				lines: &[]string{"line\n"},
			},
			tab: bufResp.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			printTab(tt.args.lines)
			got := buf.String()

			if !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("printTab() error = %v, want %v", got, tt.tab)
			}
		})
	}
}

func Test_printAliases(t *testing.T) {
	prepareToken(t)
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
			name: "printAlias test",
			args: args{
				app:   getMockApp(),
				lines: &[]string{},
			},
			lines: &[]string{"alias.name\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printAliases(tt.args.app, tt.args.lines)
			if !reflect.DeepEqual(tt.lines, tt.args.lines) {
				t.Errorf("printAliases() error = %v, want %v", tt.args.lines, tt.lines)
			}
		})
	}
}

func Test_getAlias(t *testing.T) {
	prepareToken(t)
	bufResp := bytes.NewBufferString("")
	tabWriter := tabwriter.NewWriter(bufResp, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprint(tabWriter, "NAME\n")
	fmt.Fprint(tabWriter, "alias.name\n")
	tabWriter.Flush()

	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		host    string
		args    args
		wantErr bool
		tab     string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "getAlias valid test",
			args: args{
				in0: context.Background(),
			},
			wantErr: false,
			tab:     bufResp.String(),
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
			name: "getAlias invalid test, HTTP error",
			args: args{
				in0: context.Background(),
			},
			wantErr: true,
			tab:     bufResp.String(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				rest.ERROR(w, ierrors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			cliutils.SetClient(server.URL, tt.host)
			buf := bytes.NewBufferString("")
			cliutils.SetOutput(buf)
			err := getAlias(tt.args.in0)
			got := buf.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("getAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.tab) {
				t.Errorf("getAlias() error = %v, want %v", got, tt.tab)
			}
		})
	}
}
