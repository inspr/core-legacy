package cli

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Test_methodNameByType(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := methodNameByType(tt.args.v); got != tt.want {
				t.Errorf("methodNameByType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlag_flag(t *testing.T) {
	tests := []struct {
		name string
		fl   *Flag
		want *pflag.Flag
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fl.flag(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flag.flag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reflectValueOf(t *testing.T) {
	type args struct {
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want []reflect.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reflectValueOf(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reflectValueOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFlags(t *testing.T) {
	type args struct {
		cmd   *cobra.Command
		flags []*Flag
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseFlags(tt.args.cmd, tt.args.flags)
		})
	}
}

func TestAddFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddFlags(tt.args.cmd)
		})
	}
}

func Test_hasCmdAnnotation(t *testing.T) {
	type args struct {
		cmdName     string
		annotations []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasCmdAnnotation(tt.args.cmdName, tt.args.annotations); got != tt.want {
				t.Errorf("hasCmdAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}
