package cli

import (
	"reflect"
	"testing"
)

func TestApplyChannel(t *testing.T) {
	tests := []struct {
		name string
		want RunMethod
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplyChannel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}
