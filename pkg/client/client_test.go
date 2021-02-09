package dappclient

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
	"gitlab.inspr.dev/inspr/core/pkg/sidecar/models"
)

func mockHTTPClient(addr string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(string, string) (net.Conn, error) {
				return net.Dial("unix", addr)
			},
		},
	}
}

func mockMessage() models.Message {
	return models.Message{
		Data: nil,
	}
}

func mockHandlerFunc(path string, expectedData interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(400)
			return
		}
		if r.URL.Path != path {
			w.WriteHeader(404)
			return
		}

		var body interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil && reflect.DeepEqual(body, expectedData) {
			data := ierrors.NewError().BadRequest().Build()
			rest.JSON(w, 500, data)
			return
		}
		if path == "/readMessage" {
			data := mockMessage()
			rest.JSON(w, 200, data)
			return
		}
	})
}

func mockHandlerFuncTimeout() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 5)
	})
}

func TestNewAppClient(t *testing.T) {
	tests := []struct {
		name string
		want *Client
	}{
		{
			name: "Valid App Client - with address",
			want: &Client{
				addr:  "/addr/to/socket",
				httpc: *mockHTTPClient("/addr/to/socket"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("INSPR_INPUT_CHANNELS", "inp1;inp2;inp3")
			os.Setenv("INSPR_OUTPUT_CHANNELS", "out1;out2;out3")
			os.Setenv("INSPR_UNIX_SOCKET", "/addr/to/socket")
			got := NewAppClient()
			if got.addr != tt.want.addr {
				t.Errorf("NewClient().addr = %v, want %v", got.addr, tt.want.addr)
			}
			if got.httpc.Transport == tt.want.httpc.Transport {
				t.Errorf("NewClient().httpc.Transport = %v, want %v", got.httpc.Transport, tt.want.httpc.Transport)
			}
		})
	}
}

func TestClient_WriteMessage(t *testing.T) {
	type args struct {
		channel string
		msg     models.Message
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		interruptServer bool
		cancelContext   bool
	}{
		{
			name: "Valid request",
			args: args{
				channel: "chan1",
				msg:     mockMessage(),
			},
			wantErr: false,
		},
		{
			name: "Invalid request - server died",
			args: args{
				channel: "chan1",
				msg:     models.Message{},
			},
			wantErr:         true,
			interruptServer: true,
		},
		{
			name: "Invalid request - context canceled",
			args: args{
				channel: "chan1",
				msg:     models.Message{},
			},
			wantErr:       true,
			cancelContext: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var handler http.HandlerFunc
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tt.cancelContext {
				handler = mockHandlerFuncTimeout()
			} else {
				handler = mockHandlerFunc("/writeMessage", tt.args)
			}

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			c := Client{
				addr:  s.URL,
				httpc: *http.DefaultClient,
			}

			if tt.interruptServer {
				s.Close()
			}
			if tt.cancelContext {
				go func() {
					time.Sleep(time.Second * 2)
					cancel()
				}()
			}

			err := c.WriteMessage(ctx, "chan1", mockMessage())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ReadMessage(t *testing.T) {
	type args struct {
		ctx     context.Context
		channel string
	}
	tests := []struct {
		name            string
		args            args
		want            models.Message
		wantErr         bool
		interruptServer bool
	}{
		{
			name: "Valid request",
			args: args{
				ctx:     context.Background(),
				channel: "chan1",
			},
			wantErr:         false,
			interruptServer: false,
			want:            mockMessage(),
		},
		{
			name: "Invalid request",
			args: args{
				ctx:     context.Background(),
				channel: "chan2",
			},
			wantErr:         true,
			interruptServer: true,
			want:            models.Message{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := mockHandlerFunc("/readMessage", tt.args)

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			c := Client{
				addr:  s.URL,
				httpc: *http.DefaultClient,
			}

			if tt.interruptServer {
				s.Close()
			}

			got, err := c.ReadMessage(context.Background(), "chan1")
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ReadMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ReadMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CommitMessage(t *testing.T) {
	type args struct {
		ctx     context.Context
		channel string
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		interruptServer bool
	}{
		{
			name: "Valid request",
			args: args{
				ctx:     context.Background(),
				channel: "chan1",
			},
			wantErr:         false,
			interruptServer: false,
		},
		{
			name: "Invalid request",
			args: args{
				ctx:     nil,
				channel: "chan1",
			},
			wantErr:         true,
			interruptServer: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := mockHandlerFunc("/commit", tt.args)

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			c := Client{
				addr:  s.URL,
				httpc: *http.DefaultClient,
			}

			if tt.interruptServer {
				s.Close()
			}

			err := c.CommitMessage(context.Background(), "chan1")
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CommitMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_sendRequest(t *testing.T) {
	type args struct {
		ctx     context.Context
		method  string
		reqData clientMessage
		route   string
	}
	tests := []struct {
		name            string
		args            args
		want            models.Message
		interruptServer bool
		wantErr         bool
	}{
		{
			name: "Valid request",
			args: args{
				route:  "/commit",
				ctx:    context.Background(),
				method: http.MethodPost,
				reqData: clientMessage{
					Message: models.Message{},
					Channel: "",
				},
			},
			wantErr:         false,
			want:            models.Message{},
			interruptServer: false,
		},
		{
			name: "Invalid request - invalid route",
			args: args{
				route:  "",
				ctx:    context.Background(),
				method: http.MethodPost,
				reqData: clientMessage{
					Message: models.Message{},
					Channel: "",
				},
			},
			wantErr:         true,
			want:            models.Message{},
			interruptServer: false,
		},
		{
			name: "Invalid request - server killed",
			args: args{
				route:  "/",
				ctx:    context.Background(),
				method: http.MethodPost,
				reqData: clientMessage{
					Message: models.Message{},
					Channel: "",
				},
			},
			wantErr:         true,
			want:            models.Message{},
			interruptServer: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := mockHandlerFunc(tt.args.route, tt.args)
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			c := Client{
				addr:  s.URL,
				httpc: *http.DefaultClient,
			}

			if tt.interruptServer {
				s.Close()
			}

			got, err := c.sendRequest(tt.args.ctx, tt.args.method, tt.args.route, tt.args.reqData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.sendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.sendRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
