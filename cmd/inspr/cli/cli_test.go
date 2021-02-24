package cli

import (
	"bytes"
	"testing"
)

// TestNewInsprCommand is mainly for improving test coverage,
// it was really tested by instantiating Inspr's CLI
func TestNewInsprCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Creates a new Cobra command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}
			got := NewInsprCommand(out, err)
			if got == nil {
				t.Errorf("NewInsprCommand() = %v", got)
			}
		})
	}
}
