package cli

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	cliutils "github.com/inspr/inspr/cmd/inspr/cli/utils"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"gopkg.in/yaml.v2"
)

func createSchema() string {
	schema := `{"type":"string"}`
	data, _ := yaml.Marshal(&schema)
	return string(data)
}

func TestNewApplyChannelType(t *testing.T) {
	chanTypeWithoutNameBytes, _ := yaml.Marshal(meta.ChannelType{})
	chanTypeDefaultBytes, _ := yaml.Marshal(meta.ChannelType{Meta: meta.Metadata{Name: "mock"}})
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "default_test",
			args: args{
				b: chanTypeDefaultBytes,
			},
			want: nil,
		},
		{
			name: "channeltype_without_name",
			args: args{
				b: chanTypeWithoutNameBytes,
			},
			want: ierrors.NewError().Message("channelType without name").Build(),
		},
		{
			name: "error_testing",
			args: args{
				b: chanTypeDefaultBytes,
			},
			want: errors.New("new error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cliutils.SetMockedClient(tt.want)
			got := NewApplyChannelType()

			r := got(tt.args.b, nil)

			if r != nil && tt.want != nil {
				if r.Error() != tt.want.Error() {
					t.Errorf("newApplyChannelType() = %v, want %v", r.Error(), tt.want.Error())
				}
			} else {
				if r != tt.want {
					t.Errorf("newApplyChannelType() = %v, want %v", r, tt.want)
				}
			}
		})
	}
}

func Test_schemaNeedsInjection(t *testing.T) {
	yamlString := createSchema()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		"test.schema",
		[]byte(yamlString),
		os.ModePerm,
	)
	defer os.Remove("test.schema")

	type args struct {
		schema string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Schema has path to existing file",
			args: args{
				schema: "test.schema",
			},
			want: true,
		},
		{
			name: "Schema doesn't need injection",
			args: args{
				schema: "thisisnotafilepath",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := schemaNeedsInjection(tt.args.schema); got != tt.want {
				t.Errorf("schemaNeedsInjection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_injectSchema(t *testing.T) {
	yamlString := createSchema()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		"test.schema",
		[]byte(yamlString),
		os.ModePerm,
	)
	defer os.Remove("test.schema")

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Returns valid schema",
			args: args{
				path: "test.schema",
			},
			wantErr: false,
			want:    "'{\"type\":\"string\"}'",
		},
		{
			name: "Invalid file path",
			args: args{
				path: "thisisnotafilepath",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := injectedSchema(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("injectSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotJSON interface{}
			var wantJSON interface{}
			json.Unmarshal([]byte(got), &gotJSON)
			json.Unmarshal([]byte(tt.want), &wantJSON)

			if !reflect.DeepEqual(gotJSON, wantJSON) {
				t.Errorf("injectSchema() = %v, want %v", gotJSON, wantJSON)
			}
		})
	}
}
