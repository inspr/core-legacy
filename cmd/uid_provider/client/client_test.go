package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"gitlab.inspr.dev/inspr/core/pkg/api/auth"
	"gitlab.inspr.dev/inspr/core/pkg/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
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
	setup()
	defer teardown()

	auxCtx := context.Background()
	auxUser := User{
		UID:      "user1",
		Role:     1,
		Scope:    []string{"ascope"},
		Password: "none",
	}
	auxUser2 := User{
		UID:      "user2",
		Role:     0,
		Scope:    []string{"ascope"},
		Password: "none",
	}

	strData, _ := json.Marshal(auxUser)
	redisClient.rdb.Set(auxCtx, auxUser.UID, strData, 0)
	strData2, _ := json.Marshal(auxUser2)
	redisClient.rdb.Set(auxCtx, auxUser2.UID, strData2, 0)

	type args struct {
		uid     string
		pwd     string
		newUser User
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "Creates a new user",
			c:    &redisClient,
			args: args{
				uid: auxUser.UID,
				pwd: auxUser.Password,
				newUser: User{
					UID:      "user3",
					Password: "u3pwd",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid - requestor can't create new users",
			c:    &redisClient,
			args: args{
				uid: auxUser2.UID,
				pwd: auxUser2.Password,
				newUser: User{
					UID:      "user3",
					Password: "u3pwd",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.CreateUser(auxCtx, tt.args.uid, tt.args.pwd, tt.args.newUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			createdUser, err := get(auxCtx, redisClient.rdb, tt.args.newUser.UID)
			if err != nil || reflect.DeepEqual(createdUser, User{}) {
				t.Errorf("Client.CreateUser() error = %v", err)
			}
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	setup()
	defer teardown()

	auxCtx := context.Background()
	auxUser := User{
		UID:      "user1",
		Role:     1,
		Scope:    []string{"ascope"},
		Password: "none",
	}
	auxUser2 := User{
		UID:      "user2",
		Role:     0,
		Scope:    []string{"ascope"},
		Password: "none",
	}
	auxUser3 := User{
		UID:      "user3",
		Role:     0,
		Scope:    []string{"ascope"},
		Password: "1234",
	}

	strData, _ := json.Marshal(auxUser)
	redisClient.rdb.Set(auxCtx, auxUser.UID, strData, 0)
	strData2, _ := json.Marshal(auxUser2)
	redisClient.rdb.Set(auxCtx, auxUser2.UID, strData2, 0)
	strData3, _ := json.Marshal(auxUser3)
	redisClient.rdb.Set(auxCtx, auxUser3.UID, strData3, 0)

	type args struct {
		uid            string
		pwd            string
		usrToBeDeleted string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "Deletes an existent user",
			c:    &redisClient,
			args: args{
				uid:            auxUser.UID,
				pwd:            auxUser.Password,
				usrToBeDeleted: "user3",
			},
			wantErr: false,
		},
		{
			name: "Invalid - requestor can't delete users",
			c:    &redisClient,
			args: args{
				uid:            auxUser2.UID,
				pwd:            auxUser2.Password,
				usrToBeDeleted: "user3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.DeleteUser(auxCtx, tt.args.uid, tt.args.pwd, tt.args.usrToBeDeleted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			createdUser, err := get(auxCtx, redisClient.rdb, tt.args.usrToBeDeleted)
			if err == nil || !reflect.DeepEqual(createdUser, User{}) {
				t.Errorf("Client.DeleteUser() error = %v", err)
			}
		})
	}
}

// func TestClient_UpdatePassword(t *testing.T) {
// 	type args struct {
// 		ctx            context.Context
// 		uid            string
// 		usrToBeUpdated string
// 		newPwd         string
// 	}
// 	tests := []struct {
// 		name    string
// 		c       *Client
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := tt.c.UpdatePassword(tt.args.ctx, tt.args.uid, tt.args.usrToBeUpdated, tt.args.newPwd); (err != nil) != tt.wantErr {
// 				t.Errorf("Client.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

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
	setup()
	defer teardown()

	auxCtx := context.Background()
	auxUser := User{
		UID:      "user1",
		Role:     1,
		Scope:    []string{"ascope"},
		Password: "none",
	}

	type args struct {
		data User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create new user",
			args: args{
				data: auxUser,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := set(auxCtx, redisClient.rdb, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			createdUser, err := get(auxCtx, redisClient.rdb, auxUser.UID)
			if err != nil || reflect.DeepEqual(createdUser, User{}) {
				t.Errorf("user wasn't set. Error %v", err)
			}
		})
	}
}

func Test_get(t *testing.T) {
	setup()
	defer teardown()

	auxCtx := context.Background()
	auxUser := User{
		UID:      "user1",
		Role:     1,
		Scope:    []string{"ascope"},
		Password: "none",
	}

	strData, _ := json.Marshal(auxUser)
	redisClient.rdb.Set(auxCtx, auxUser.UID, strData, 0)

	type args struct {
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
		{
			name: "Invalid - get user given non-existent UID",
			args: args{
				key: "RANDOMKEY",
			},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := get(auxCtx, redisClient.rdb, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				fmt.Println(got)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_delete(t *testing.T) {
	setup()
	defer teardown()

	auxCtx := context.Background()
	auxUser := User{
		UID:      "user1",
		Role:     1,
		Scope:    []string{"ascope"},
		Password: "none",
	}

	strData, _ := json.Marshal(auxUser)
	redisClient.rdb.Set(auxCtx, auxUser.UID, strData, 0)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete user given UID",
			args: args{
				key: "user1",
			},
			wantErr: false,
		},
		{
			name: "Invalid - Delete non-existent user",
			args: args{
				key: "RANDOMUSER",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := delete(auxCtx, redisClient.rdb, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			deletedUser, err := get(auxCtx, redisClient.rdb, tt.args.key)
			if err == nil {
				t.Errorf("user %v wasn't deleted", deletedUser)
			}
		})
	}
}

func Test_havePermission(t *testing.T) {
	setup()
	defer teardown()

	auxCtx := context.Background()
	auxUser := User{
		UID:      "user1",
		Role:     1,
		Scope:    []string{"ascope"},
		Password: "none",
	}
	auxUser2 := User{
		UID:      "user2",
		Role:     0,
		Scope:    []string{"ascope"},
		Password: "none",
	}

	strData, _ := json.Marshal(auxUser)
	redisClient.rdb.Set(auxCtx, auxUser.UID, strData, 0)
	strData2, _ := json.Marshal(auxUser2)
	redisClient.rdb.Set(auxCtx, auxUser2.UID, strData2, 0)

	type args struct {
		uid string
		pwd string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "User has admin permisson",
			args: args{
				uid: "user1",
				pwd: "none",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Non existent user",
			args: args{
				uid: "RANDOMUSER",
				pwd: "",
			},
			wantErr: true,
		},
		{
			name: "Wrong user credentials",
			args: args{
				uid: "user1",
				pwd: "invalid",
			},
			wantErr: true,
		},
		{
			name: "User is not admin",
			args: args{
				uid: "user2",
				pwd: "none",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := havePermission(auxCtx, redisClient.rdb, tt.args.uid, tt.args.pwd)

			if (err != nil) != tt.wantErr {
				t.Errorf("delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("havePermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encrypt(t *testing.T) {
	os.Setenv("REFRESH_URL", "randomurl")
	os.Setenv("REFRESH_KEY", "61626364616263646162636461626364")
	defer os.Unsetenv("REFRESH_KEY")
	defer os.Unsetenv("REFRESH_URL")

	type args struct {
		user User
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Returns payload with encrypted refresh token",
			args: args{
				user: User{
					UID:      "user1",
					Role:     1,
					Scope:    []string{"ascope"},
					Password: "none",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encrypt(tt.args.user)
			if err != nil {
				t.Errorf("encrypt() return an error: %v", err)
				return
			}

			if got.Refresh == "" {
				t.Errorf("error while creating the refresh token")
				return
			}
		})
	}
}

func Test_decrypt(t *testing.T) {
	os.Setenv("REFRESH_URL", "randomurl")
	os.Setenv("REFRESH_KEY", "61626364616263646162636461626364")
	defer os.Unsetenv("REFRESH_KEY")
	defer os.Unsetenv("REFRESH_URL")

	auxUser := User{
		UID:      "user1",
		Password: "strongpwd",
	}

	payload, _ := encrypt(auxUser)

	tests := []struct {
		name string
		want User
	}{
		{
			name: "Decrypts valid user refresh token",
			want: User{
				UID:      "user1",
				Password: "strongpwd",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decrypt(payload.Refresh)
			if err != nil {
				t.Errorf("unable to decrypt, error: %v", err)
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
		hand    func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
				payload: auth.Payload{
					UID:   "user1",
					Scope: []string{"app1", "app2"},
				},
			},
			want:    "user1-app1-app2",
			wantErr: false,
			hand: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					rest.ERROR(w, fmt.Errorf("method should be POST"))
				}
				data := auth.Payload{}
				if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
					rest.ERROR(w, err)
					return
				}
				strScope := strings.Join(data.Scope, "-")
				token := fmt.Sprintf("%s-%s", data.UID, strScope)
				val := models.AuthDI{
					Token: token,
				}
				rest.JSON(w, 200, val)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.hand))
			defer server.Close()
			os.Setenv("INSPR_CLUSTER_ADDR", server.URL)
			defer os.Unsetenv("INSPR_CLUSTER_ADDR")

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
