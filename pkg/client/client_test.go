package dappclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/rest/request"
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.URL.Path != path {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var body interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil && reflect.DeepEqual(body, expectedData) {
			data := ierrors.New("").BadRequest()
			rest.JSON(w, http.StatusInternalServerError, data)
			return
		}
		if path == "/readMessage" {
			data := mockMessage()
			rest.JSON(w, http.StatusOK, data)
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
	os.Setenv("INSPR_APP_SCOPE", "random.ctx")
	os.Setenv("INSPR_ENV", "test")
	os.Setenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS", "kafka")
	os.Setenv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET", "latest")
	os.Setenv("ch1_SCHEMA", `{"type":"string"}`)
	os.Setenv("ch2_SCHEMA", "hellotest")
	os.Setenv("INSPR_APP_ID", "testappid1")
	os.Setenv("INSPR_LBSIDECAR_IMAGE", "random-sidecar-image")
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockEnvVars() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_APP_SCOPE")
	os.Unsetenv("INSPR_ENV")
	os.Unsetenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS")
	os.Unsetenv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET")
	os.Unsetenv("INSPR_APP_ID")
	os.Unsetenv("INSPR_LBSIDECAR_IMAGE")
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
				handler = mockHandlerFunc("/channel/chan1", tt.args)
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
			err := client.Send(
				context.Background(),
				"/channel/"+tt.args.channel,
				http.MethodPost,
				struct{ Message interface{} }{tt.message},
				&response)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client_HandleChannel response.Status = %v, wantErr = %v", response.Status, tt.wantErr)
			}
		})
	}
}

func TestClient_HandleRoute(t *testing.T) {
	type args struct {
		path    string
		handler func(t *testing.T) func(w http.ResponseWriter, r *http.Request)
	}
	tests := []struct {
		name    string
		args    args
		msg     string
		wantErr bool
	}{
		{
			name: "valid route handle",
			msg:  "message",
			args: args{
				path: "hello/world",
				handler: func(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
					return func(w http.ResponseWriter, r *http.Request) {
						var msg string
						decoder := json.NewDecoder(r.Body)
						err := decoder.Decode(&msg)
						if err != nil {
							t.Errorf("Client_HandleRoute message error = %v", err)
						}

						if msg != "message" {
							fmt.Println(msg)
							t.Errorf("Client_HandleRoute message = %v, want message", msg)
						}

						rest.JSON(w, http.StatusOK, nil)
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
			c.HandleRoute(tt.args.path, tt.args.handler(t))
			s := httptest.NewServer(c.mux)
			client := request.NewJSONClient(s.URL)
			response := struct {
				Status string `json:"status"`
			}{}
			err := client.Send(
				context.Background(),
				"/route/"+tt.args.path,
				http.MethodPost,
				tt.msg,
				&response)

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
				readAddr: ":3304",
			}
			client.HandleChannel(tt.channel, handler)
			errch := make(chan error)
			go func() {
				errch <- client.Run(ctx)
			}()

			c := request.NewJSONClient("http://localhost:3304/channel")
			var response struct {
				Status string `json:"status"`
			}

			err := c.Send(
				ctx,
				tt.channel,
				http.MethodPost,
				tt.message,
				&response)

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

func TestClient_SendRequest(t *testing.T) {
	type args struct {
		nodeName string
		path     string
		method   string
		body     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			name: "valid route request",
			args: args{
				nodeName: "rt1",
				path:     "end1",
				method:   "POST",
				body:     "test_body",
			},
			wantErr: false,
			want:    "hello",
		},
		{
			name: "Invalid route request",
			args: args{
				nodeName: "rt1",
				path:     "end1",
				method:   "GET",
				body:     "test_body",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSever := createMockedLBsidecar("5100", tt.args.nodeName, tt.args.path, tt.args.method, tt.want, tt.args.body)

			testSever.Start()
			defer testSever.Close()

			c := Client{
				client: request.NewJSONClient(testSever.URL).
					HTTPClient(*http.DefaultClient).
					Pointer(),
			}

			var resp string
			err := c.SendRequest(context.Background(), tt.args.nodeName, tt.args.path, tt.args.method, tt.args.body, &resp)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Client.SendRequest() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if resp != tt.want {
				t.Errorf("Client.SendRequest() = %v, want %v", resp, tt.want)
				return
			}
		})
	}
}

func createMockedLBsidecar(port, nodeName, path, method, want string, body interface{}) *httptest.Server {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var receivedData string
			err := json.NewDecoder(r.Body).Decode(&receivedData)
			if err != nil {
				rest.ERROR(w, fmt.Errorf("error while reading msg body"))
				return
			}

			if receivedData != body {
				rest.ERROR(w, fmt.Errorf("invalid messagem body"))
				return
			}

			if r.URL.Path != fmt.Sprintf("/route/%s/%s", nodeName, path) {
				rest.ERROR(w, fmt.Errorf("invalid path"))
				return
			}

			if r.Method != method {
				rest.ERROR(w, fmt.Errorf("invalid method"))
				return
			}

			if want == "" {
				rest.ERROR(w, fmt.Errorf("bad request"))
				return
			}

			rest.JSON(w, http.StatusOK, want)
		}),
	)

	mockServer.Listener.Close()
	mockServer.Listener = listener

	return mockServer

}
