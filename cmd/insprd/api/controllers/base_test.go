package controller

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/fake"
)

func TestServer_Init(t *testing.T) {
	type args struct {
		mm memory.Manager
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		{
			name: "successful_server_init",
			s:    &Server{},
			args: args{
				mm: fake.MockMemoryManager(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Init(tt.args.mm)
			if !reflect.DeepEqual(tt.s.MemoryManager, fake.MockMemoryManager(nil)) {
				t.Errorf("TestServer_Init() = %v, want %v", tt.s.MemoryManager, nil)
			}
		})
	}
}
