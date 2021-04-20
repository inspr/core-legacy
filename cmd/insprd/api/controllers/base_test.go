package controller

import (
	"reflect"
	"testing"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/memory/fake"
	"github.com/inspr/inspr/cmd/insprd/operators"
	ofake "github.com/inspr/inspr/cmd/insprd/operators/fake"
	"github.com/inspr/inspr/pkg/auth"
	authmock "github.com/inspr/inspr/pkg/auth/mocks"
)

func TestServer_Init(t *testing.T) {
	type args struct {
		mm   memory.Manager
		op   operators.OperatorInterface
		auth auth.Auth
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
				mm:   fake.MockMemoryManager(nil),
				op:   ofake.NewFakeOperator(),
				auth: authmock.NewMockAuth(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Init(tt.args.mm, tt.args.op, tt.args.auth)
			if !reflect.DeepEqual(tt.s.MemoryManager, fake.MockMemoryManager(nil)) {
				t.Errorf("TestServer_Init() = %v, want %v", tt.s.MemoryManager, nil)
			}
		})
	}
}
