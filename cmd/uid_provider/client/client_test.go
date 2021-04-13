package client

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"gitlab.inspr.dev/inspr/core/pkg/api/auth"
)

var redisServer *miniredis.Miniredis
var redisClient Client

func TestNewRedisClient(t *testing.T) {
	tests := []struct {
		name string
		want *Client
	}{
		{
			name: "client_creation",
			want: &Client{
				rdb: &redis.Client{},
			},
		},
	}
	for _, tt := range tests {
		got := NewRedisClient()

		if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
			t.Errorf(
				"NewRedisClient() = %v, want %v",
				reflect.TypeOf(got),
				reflect.TypeOf(tt.want),
			)
		}
	}
}

func TestClient_CreateUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		uid     string
		newUser User
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.CreateUser(tt.args.ctx, tt.args.uid, tt.args.newUser); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	type args struct {
		ctx            context.Context
		uid            string
		usrToBeDeleted string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.DeleteUser(tt.args.ctx, tt.args.uid, tt.args.usrToBeDeleted); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UpdatePassword(t *testing.T) {
	type args struct {
		ctx            context.Context
		uid            string
		usrToBeUpdated string
		newPwd         string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UpdatePassword(tt.args.ctx, tt.args.uid, tt.args.usrToBeUpdated, tt.args.newPwd); (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		uid string
		pwd string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Login(tt.args.ctx, tt.args.uid, tt.args.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RefreshToken(t *testing.T) {
	type args struct {
		ctx          context.Context
		refreshToken string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		want    auth.Payload
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.RefreshToken(tt.args.ctx, tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.RefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set(t *testing.T) {
	type args struct {
		ctx  context.Context
		rdb  *redis.Client
		data User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := set(tt.args.ctx, tt.args.rdb, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_get(t *testing.T) {
	setup()
	defer teardown()

	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr bool
	}{
		{
			name: "Get user given UID",
			args: args{
				ctx: context.Background(),
				key: "user1",
			},
			want: User{
				UID:      "user1",
				Role:     1,
				Scope:    []string{"ascope"},
				Password: "none",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strData, _ := json.Marshal(tt.want)
			redisClient.rdb.Set(tt.args.ctx, tt.args.key, strData, 0)

			got, err := get(tt.args.ctx, redisClient.rdb, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_delete(t *testing.T) {
	type args struct {
		ctx context.Context
		rdb *redis.Client
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := delete(tt.args.ctx, tt.args.rdb, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isAdmin(t *testing.T) {
	type args struct {
		ctx context.Context
		rdb *redis.Client
		uid string
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
			if got := isAdmin(tt.args.ctx, tt.args.rdb, tt.args.uid); got != tt.want {
				t.Errorf("isAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encrypt(t *testing.T) {
	type args struct {
		user User
	}
	tests := []struct {
		name    string
		args    args
		want    auth.Payload
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encrypt(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decrypt(t *testing.T) {
	type args struct {
		encryptedString string
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decrypt(tt.args.encryptedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_requestNewToken(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload auth.Payload
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requestNewToken(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("requestNewToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("requestNewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Auxiliar methods

func setup() {
	redisServer = mockRedis()
	redisClient.rdb = redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
}

func teardown() {
	redisServer.Close()
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}
