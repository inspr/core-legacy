package controller

import (
	"reflect"
	"testing"

	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/fake"
	"inspr.dev/inspr/cmd/insprd/operators"
	ofake "inspr.dev/inspr/cmd/insprd/operators/fake"
	"inspr.dev/inspr/pkg/auth"
	authmock "inspr.dev/inspr/pkg/auth/mocks"
)

func TestServer_Init(t *testing.T) {
	type args struct {
		mem  memory.Manager
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
				mem:  fake.GetMockMemoryManager(nil, nil),
				op:   ofake.NewFakeOperator(),
				auth: authmock.NewMockAuth(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Init(tt.args.mem, tt.args.op, tt.args.auth)
			fake := fake.MockTreeMemory(nil)
			tree := tt.s.memory.Tree()
			if !reflect.DeepEqual(tree, fake) {
				t.Errorf("TestServer_Init() = %v, want %v", tree, fake)
			}
		})
	}
}
