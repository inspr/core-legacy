package cli

import (
	"errors"
	"testing"

	cliutils "github.com/inspr/inspr/cmd/inspr/cli/utils"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"gopkg.in/yaml.v2"
)

func TestNewApplyAlias(t *testing.T) {
	chanWithoutNameBytes, _ := yaml.Marshal(meta.Alias{})
	chanDefaultBytes, _ := yaml.Marshal(meta.Alias{Meta: meta.Metadata{Name: "mock"}})
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
				b: chanDefaultBytes,
			},
			want: nil,
		},
		{
			name: "alias_without_name",
			args: args{
				b: chanWithoutNameBytes,
			},
			want: ierrors.NewError().Message("alias without name").Build(),
		},
		{
			name: "error_testing",
			args: args{
				b: chanDefaultBytes,
			},
			want: errors.New("new error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cliutils.SetMockedClient(tt.want)
			got := NewApplyAlias()

			r := got(tt.args.b, nil)

			if r != nil && tt.want != nil {
				if r.Error() != tt.want.Error() {
					t.Errorf("NewApplyAlias() = %v, want %v", r.Error(), tt.want.Error())
				}
			} else {
				if r != tt.want {
					t.Errorf("NewApplyAlias() = %v, want %v", r, tt.want)
				}
			}
		})
	}
}
