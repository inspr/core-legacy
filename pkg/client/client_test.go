package dappclient

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/rest/request"
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

func mockMessage() interface{} {
	return nil
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

// createMockEnvVars - sets up the env values to be used in the tests functions
func createMockEnvVars() {
	os.Setenv("INSPR_INPUT_CHANNELS", "inp1;inp2;inp3")
	os.Setenv("INSPR_OUTPUT_CHANNELS", "inp1;inp2;inp3")
	os.Setenv("INSPR_UNIX_SOCKET", "/addr/to/socket")
	os.Setenv("INSPR_APP_CTX", "random.ctx")
	os.Setenv("INSPR_ENV", "test")
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "kafka")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "latest")
	os.Setenv("ch1_SCHEMA", `{"type":"string"}`)
	os.Setenv("ch2_SCHEMA", "hellotest")
	os.Setenv("INSPR_APP_ID", "testappid1")
	os.Setenv("INSPR_SIDECAR_IMAGE", "random-sidecar-image")
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockEnvVars() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_APP_CTX")
	os.Unsetenv("INSPR_ENV")
	os.Unsetenv("KAFKA_BOOTSTRAP_SERVERS")
	os.Unsetenv("KAFKA_AUTO_OFFSET_RESET")
	os.Unsetenv("INSPR_APP_ID")
	os.Unsetenv("INSPR_SIDECAR_IMAGE")
}

func TestNewAppClient(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()
	tests := []struct {
		name string
		want *Client
	}{
		{
			name: "Valid App Client - with address",
			want: &Client{
				client: request.NewClient().
					BaseURL("http://unix").
					HTTPClient(*mockHTTPClient("http://unix")).
					Encoder(json.Marshal).
					Decoder(request.JSONDecoderGenerator).
					Pointer(),
			},
		},
	}
	environment.SetMockEnv()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAppClient()
			if got.client == nil {
				t.Errorf("NewClient().client = %v, want %v", got.client, tt.want.client)
			}
		})
	}
}

func TestClient_WriteMessage(t *testing.T) {
	type args struct {
		channel string
		msg     interface{}
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
				msg:     nil,
			},
			wantErr:         true,
			interruptServer: true,
		},
		{
			name: "Invalid request - context canceled",
			args: args{
				channel: "chan1",
				msg:     nil,
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
				handler = mockHandlerFunc("/chan1", tt.args)
			}

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			c := Client{
				client: request.NewClient().
					BaseURL(s.URL).
					HTTPClient(*http.DefaultClient).
					Encoder(json.Marshal).
					Decoder(request.JSONDecoderGenerator).
					Pointer(),
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

func TestClient_HandleChannel(t *testing.T) {
	type fields struct {
		readAddr string
	}
	type args struct {
		channel string
		handler func(t *testing.T) func(ctx context.Context, body io.Reader) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		message interface{}
	}{
		{
			name:    "no error on handler",
			message: "message",
			args: args{
				channel: "channel",
				handler: func(t *testing.T) func(ctx context.Context, body io.Reader) error {
					return func(ctx context.Context, body io.Reader) error {

						message := struct {
							Message string
						}{}
						decoder := json.NewDecoder(body)
						err := decoder.Decode(&message)
						if err != nil {
							return err
						}
						if message.Message != "message" {
							t.Errorf("Client_HandleChannel message = %v, want message", message.Message)
						}

						return nil
					}
				},
			},
		},
		{
			name:    "error on handler",
			message: "message",
			wantErr: true,
			args: args{
				channel: "channel",
				handler: func(t *testing.T) func(ctx context.Context, body io.Reader) error {
					return func(ctx context.Context, body io.Reader) error {

						message := struct {
							Message string
						}{}
						decoder := json.NewDecoder(body)
						err := decoder.Decode(&message)
						if err != nil {
							return err
						}
						if message.Message != "message" {
							t.Errorf("Client_HandleChannel message = %v, want message", message.Message)
						}

						return errors.New("Error")
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				mux: http.NewServeMux(),
			}
			c.HandleChannel(tt.args.channel, tt.args.handler(t))
			s := httptest.NewServer(c.mux)
			client := request.NewJSONClient(s.URL)
			response := struct {
				Status string `json:"status"`
			}{}
			err := client.Send(context.Background(), tt.args.channel, "POST", struct{ Message interface{} }{tt.message}, &response)
			if (response.Status != "OK") != tt.wantErr {
				t.Errorf("Client_HandleChannel response.Status = %v, wantErr = %v", response.Status, tt.wantErr)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Client_HandleChannel response.Status = %v, wantErr = %v", response.Status, tt.wantErr)
			}
		})
	}
}

func TestClient_Run(t *testing.T) {
	environment.SetMockEnv()
	defer environment.UnsetMockEnv()
	tests := []struct {
		name    string
		message interface{}
		wantErr bool
		handler func(t *testing.T) func(ctx context.Context, r io.Reader) error
		channel string
	}{
		{
			message: "this is a message",
			name:    "correct functionality",
			channel: "banana",
		},

		{
			message: "this is a message",
			name:    "incorrect functionality",
			channel: "banana",
			wantErr: true,
			handler: func(t *testing.T) func(ctx context.Context, r io.Reader) error {
				return func(ctx context.Context, r io.Reader) error {
					return errors.New("this is an error")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var handler func(context.Context, io.Reader) error
			if tt.handler != nil {
				handler = tt.handler(t)
			} else {
				handler = func(ctx context.Context, r io.Reader) error {
					decoder := json.NewDecoder(r)
					var response interface{}
					err := decoder.Decode(&response)
					if err != nil {
						return err
					}
					if !reflect.DeepEqual(response, tt.message) {
						t.Errorf("Client_Run response = %v, want %v", response, tt.message)
					}
					return nil
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			client := &Client{
				mux:      http.NewServeMux(),
				readAddr: ":3301",
			}
			client.HandleChannel(tt.channel, handler)
			errch := make(chan error)
			go func() {
				errch <- client.Run(ctx)
			}()

			c := request.NewJSONClient("http://localhost:3301")
			var response struct {
				Status string `json:"status"`
			}
			err := c.Send(ctx, tt.channel, "POST", tt.message, &response)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client_Run err = %v, wantErr = %v", err, tt.wantErr)
			}
			cancel()
			err = <-errch
			if err != nil && err != context.Canceled {
				t.Errorf("Client_Run error in server = %v, wantErr = %v", err, tt.wantErr)
			}

		})
	}
}
