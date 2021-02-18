package cli

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewInsprCommand(t *testing.T) {
	tests := []struct {
		name    string
		want    *cobra.Command
		wantOut string
		wantErr string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}
			if got := NewInsprCommand(out, err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInsprCommand() = %v, want %v", got, tt.want)
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("NewInsprCommand() = %v, want %v", gotOut, tt.wantOut)
			}
			if gotErr := err.String(); gotErr != tt.wantErr {
				t.Errorf("NewInsprCommand() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
