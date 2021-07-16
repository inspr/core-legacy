package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"inspr.dev/inspr/cmd/uid_provider/api/models"
	"inspr.dev/inspr/cmd/uid_provider/client"
)

var redisServer *miniredis.Miniredis
var redisClient client.Client
var insprServer *httptest.Server

func TestNewHandler(t *testing.T) {
	setup()
	defer teardown()
	auxCtx := context.Background()

	type args struct {
		rdb client.RedisManager
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Returns a new handler",
			args: args{
				client.NewRedisClient(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHandler(auxCtx, tt.args.rdb)
			if got == nil {
				t.Errorf("NewHandler() error: %v", got)
			}
		})
	}
}

func TestHandler_CreateUserHandler(t *testing.T) {
	setup()
	defer teardown()
	auxCtx := context.Background()

	tests := []struct {
		name string
		h    *Handler
		body models.ReceivedDataCreate
		want int
	}{
		{
			name: "Send request to CreateUserHandler",
			h:    NewHandler(auxCtx, &redisClient),
			body: models.ReceivedDataCreate{
				UID:      "rand",
				Password: "123",
				User:     client.User{},
			},
			want: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.h.CreateUserHandler())
			defer ts.Close()

			client := ts.Client()
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Log("error decoding payload into bytes")
				return
			}
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()
			if res.StatusCode != tt.want {
				t.Errorf("CreateUserHandler() = %v, want %v", res.StatusCode, tt.want)
				return
			}
		})
	}
}

func TestHandler_DeleteUserHandler(t *testing.T) {
	setup()
	defer teardown()
	auxCtx := context.Background()

	tests := []struct {
		name string
		h    *Handler
		body models.ReceivedDataDelete
		want int
	}{
		{
			name: "Send request to DeleteUserHandler",
			h:    NewHandler(auxCtx, &redisClient),
			body: models.ReceivedDataDelete{
				UID:             "rand",
				Password:        "123",
				UserToBeDeleted: "rand2",
			},
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.h.DeleteUserHandler())
			defer ts.Close()

			client := ts.Client()
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Log("error decoding payload into bytes")
				return
			}
			req, _ := http.NewRequest(http.MethodDelete, ts.URL, bytes.NewBuffer(body))
			res, err := client.Do(req)
			if err != nil {
				t.Log("error making a PUT in the httptest server")
				return
			}
			defer res.Body.Close()
			if res.StatusCode != tt.want {
				t.Errorf("DeleteUserHandler() = %v, want %v", res.StatusCode, tt.want)
				return
			}
		})
	}
}

func TestHandler_UpdatePasswordHandler(t *testing.T) {
	setup()
	defer teardown()
	auxCtx := context.Background()

	tests := []struct {
		name string
		h    *Handler
		body models.ReceivedDataUpdate
		want int
	}{
		{
			name: "Send request to UpdatePasswordHandler",
			h:    NewHandler(auxCtx, &redisClient),
			body: models.ReceivedDataUpdate{
				UID:             "rand",
				Password:        "123",
				UserToBeUpdated: "rand2",
				NewPassword:     "321",
			},
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.h.UpdatePasswordHandler())
			defer ts.Close()

			client := ts.Client()
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Log("error decoding payload into bytes")
				return
			}
			req, _ := http.NewRequest(http.MethodPut, ts.URL, bytes.NewBuffer(body))
			res, err := client.Do(req)
			if err != nil {
				t.Log("error making a PUT in the httptest server")
				return
			}
			defer res.Body.Close()
			if res.StatusCode != tt.want {
				t.Errorf("UpdatePasswordHandler() = %v, want %v", res.StatusCode, tt.want)
				return
			}
		})
	}
}

func TestHandler_LoginHandler(t *testing.T) {
	setup()
	defer teardown()
	auxCtx := context.Background()

	tests := []struct {
		name string
		h    *Handler
		body models.ReceivedDataLogin
		want int
	}{
		{
			name: "Send request to LoginHandler",
			h:    NewHandler(auxCtx, &redisClient),
			body: models.ReceivedDataLogin{
				UID:      "rand",
				Password: "123",
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.h.LoginHandler())
			defer ts.Close()

			client := ts.Client()
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Log("error decoding payload into bytes")
				return
			}
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()
			if res.StatusCode != tt.want {
				t.Errorf("LoginHandler() = %v, want %v", res.StatusCode, tt.want)
				return
			}
		})
	}
}

func TestHandler_RefreshTokenHandler(t *testing.T) {
	setup()
	defer teardown()
	auxCtx := context.Background()

	tests := []struct {
		name string
		h    *Handler
		body models.ReceivedDataRefresh
		want int
	}{
		{
			name: "Send request to RefreshTokenHandler",
			h:    NewHandler(auxCtx, &redisClient),
			body: models.ReceivedDataRefresh{
				RefreshToken: []byte("rand"),
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.h.RefreshTokenHandler())
			defer ts.Close()

			client := ts.Client()
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Log("error decoding payload into bytes")
				return
			}
			res, err := client.Post(ts.URL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Log("error making a POST in the httptest server")
				return
			}
			defer res.Body.Close()
			if res.StatusCode != tt.want {
				t.Errorf("RefreshTokenHandler() = %v, want %v", res.StatusCode, tt.want)
				return
			}
		})
	}
}

func setup() {
	redisServer = mockRedis()
	insprServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	os.Setenv("INSPR_CLUSTER_ADDR", insprServer.URL)
	os.Setenv("REFRESH_URL", "randomurl")
	os.Setenv("REFRESH_KEY", "61626364616263646162636461626364")
	os.Setenv("REDIS_HOST", redisServer.Host())
	os.Setenv("REDIS_PORT", redisServer.Port())
	os.Setenv("REDIS_PASSWORD", "")

	redisClient = *client.NewRedisClient()
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
