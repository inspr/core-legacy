package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/inspr/inspr/cmd/uid_provider/client"
)

var redisServer *miniredis.Miniredis
var redisClient client.Client
var insprServer *httptest.Server

func TestServer_Init(t *testing.T) {
	setup()
	defer teardown()
	auxCtx := context.Background()

	type args struct {
		rdb client.RedisManager
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		{
			name: "Initializes server",
			s:    &Server{},
			args: args{
				client.NewRedisClient(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Init(tt.args.rdb, auxCtx)
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
