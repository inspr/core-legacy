package controller

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/mocks"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
)

func TestServer_Init(t *testing.T) {
	type args struct {
		mm  memory.Manager
		cOp operators.ChannelOperatorInterface
		nOp operators.NodeOperatorInterface
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
				mm: mocks.MockMemoryManager(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Init(tt.args.mm, tt.args.nOp, tt.args.cOp)
			if !reflect.DeepEqual(tt.s.MemoryManager, mocks.MockMemoryManager(nil)) {
				t.Errorf("TestServer_Init() = %v, want %v", tt.s.MemoryManager, nil)
			}
		})
	}
}
