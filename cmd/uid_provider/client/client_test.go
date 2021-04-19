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
	"github.com/inspr/inspr/pkg/api/auth"
	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/rest"
)

var redisServer *miniredis.Miniredis
var redisClient Client
var insprServer *httptest.Server

func TestNewRedisClient(t *testing.T) {
	setup()
	defer teardown()
	tests := []struct {
		name string
		want *Client
	}{
		{
			name: "client_creation",
			want: &Client{
				rdb: &redis.ClusterClient{},
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
			if !tt.wantErr && (err == nil || createdUser != nil) {
				t.Errorf("Client.DeleteUser() error = %v", err)
			}
		})
	}
}

func TestClient_UpdatePassword(t *testing.T) {
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
		uid            string
		pwd            string
		usrToBeUpdated string
		newPwd         string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "Updated an user password",
			c:    &redisClient,
			args: args{
				uid:            auxUser.UID,
				pwd:            auxUser.Password,
				usrToBeUpdated: auxUser2.UID,
				newPwd:         "banana",
			},
			wantErr: false,
		},
		{
			name: "Invalid - requestor can't update users",
			c:    &redisClient,
			args: args{
				uid:            auxUser2.UID,
				pwd:            auxUser2.Password,
				usrToBeUpdated: auxUser2.UID,
				newPwd:         "banana",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.UpdatePassword(auxCtx, tt.args.uid, tt.args.pwd, tt.args.usrToBeUpdated, tt.args.newPwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			updatedUser, err := get(auxCtx, redisClient.rdb, tt.args.usrToBeUpdated)
			if !tt.wantErr && (err != nil || updatedUser.Password != tt.args.newPwd) {
				t.Errorf("Client.UpdatePassword() error = %v", err)
			}
		})
	}
}

func TestClient_Login(t *testing.T) {
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
		{
			name: "Login with a valid user",
			c:    &redisClient,
			args: args{
				uid: auxUser.UID,
				pwd: auxUser.Password,
			},
			want:    "user1-ascope",
			wantErr: false,
		},
		{
			name: "Invalid - non existent user",
			c:    &redisClient,
			args: args{
				uid: "RANDOMUSER",
				pwd: "RANDOMPWD",
			},
			wantErr: true,
		},
		{
			name: "Invalid - invalid password",
			c:    &redisClient,
			args: args{
				uid: auxUser.UID,
				pwd: "RANDOMPWD",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Login(auxCtx, tt.args.uid, tt.args.pwd)
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
		Role:     1,
		Scope:    []string{"ascope"},
		Password: "none",
	}

	strData, _ := json.Marshal(auxUser)
	redisClient.rdb.Set(auxCtx, auxUser.UID, strData, 0)

	payload, _ := redisClient.encrypt(auxUser)
	payload2, _ := redisClient.encrypt(auxUser2)

	type args struct {
		refreshToken []byte
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "Refreshing a valid token",
			c:    &redisClient,
			args: args{
				payload.Refresh,
			},
			wantErr: false,
		},
		{
			name: "Invalid - invalid token",
			c:    &redisClient,
			args: args{
				[]byte("invalid"),
			},
			wantErr: true,
		},
		{
			name: "Invalid - user was deleted",
			c:    &redisClient,
			args: args{
				payload2.Refresh,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.RefreshToken(auxCtx, tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("Client.RefreshToken() return is empty")
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
		want    *User
		wantErr bool
	}{
		{
			name: "Get user given UID",
			args: args{
				key: "user1",
			},
			want: &User{
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
			want:    nil,
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

func Test_hasPermission(t *testing.T) {
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
			err := hasPermission(auxCtx, redisClient.rdb, tt.args.uid, tt.args.pwd)

			if (err != nil) != tt.wantErr {
				t.Errorf("delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_encrypt(t *testing.T) {
	setup()
	defer teardown()

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
			got, err := redisClient.encrypt(tt.args.user)
			if err != nil {
				t.Errorf("encrypt() return an error: %v", err)
				return
			}

			if string(got.Refresh) == "" {
				t.Errorf("error while creating the refresh token")
				return
			}
		})
	}
}

func Test_decrypt(t *testing.T) {
	setup()
	defer teardown()

	auxUser := User{
		UID:      "user1",
		Password: "strongpwd",
	}

	payload, _ := redisClient.encrypt(auxUser)

	tests := []struct {
		name string
		want *User
	}{
		{
			name: "Decrypts valid user refresh token",
			want: &User{
				UID:      "user1",
				Password: "strongpwd",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := redisClient.decrypt(payload.Refresh)
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
	setup()
	defer teardown()

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
		{
			name: "Sends request for new token",
			args: args{
				ctx: context.Background(),
				payload: auth.Payload{
					UID:   "user1",
					Scope: []string{"app1", "app2"},
				},
			},
			want:    "user1-app1-app2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := redisClient.requestNewToken(tt.args.ctx, tt.args.payload)
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
	insprServer = httptest.NewServer(http.HandlerFunc(insprServerHandler))

	os.Setenv("INSPR_CLUSTER_ADDR", insprServer.URL)
	os.Setenv("REFRESH_URL", "randomurl")
	os.Setenv("REFRESH_KEY", "61626364616263646162636461626364")
	os.Setenv("REDIS_HOST", redisServer.Host())
	os.Setenv("REDIS_PORT", redisServer.Port())
	os.Setenv("REDIS_PASSWORD", "")

	redisClient = *NewRedisClient()
}

func teardown() {
	os.Unsetenv("REFRESH_KEY")
	os.Unsetenv("REFRESH_URL")
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("INSPR_CLUSTER_ADDR")
	redisServer.Close()
	insprServer.Close()
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	return s
}

func insprServerHandler(w http.ResponseWriter, r *http.Request) {
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
}
