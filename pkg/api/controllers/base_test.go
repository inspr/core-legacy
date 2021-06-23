package controller

import (
	"reflect"
	"testing"

	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/memory/fake"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	"inspr.dev/inspr/cmd/insprd/operators"
	ofake "inspr.dev/inspr/cmd/insprd/operators/fake"
	"inspr.dev/inspr/pkg/auth"
	authmock "inspr.dev/inspr/pkg/auth/mocks"
)

func TestServer_Init(t *testing.T) {
	type args struct {
		mm      tree.Manager
		op      operators.OperatorInterface
		auth    auth.Auth
		brokers brokers.Manager
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
				mm:      fake.MockMemoryManager(nil),
				op:      ofake.NewFakeOperator(),
				auth:    authmock.NewMockAuth(nil),
				brokers: fake.MockBrokerManager(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Init(tt.args.mm, tt.args.op, tt.args.auth, tt.args.brokers)
			if !reflect.DeepEqual(tt.s.MemoryManager, fake.MockMemoryManager(nil)) {
				t.Errorf("TestServer_Init() = %v, want %v", tt.s.MemoryManager, nil)
			}
		})
	}
}
