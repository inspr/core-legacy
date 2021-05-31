package sidecarserv

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/inspr/inspr/pkg/sidecar_old/models"
)

// createMockEnvVars - sets up the env values to be used in the tests functions
func createMockEnvVars() {
	customEnvValues := "chan_someBroker;testing_someBroker;banana_someBroker"
	var unixSocketAddr = "/tmp/insprd.sock"
	os.Setenv("INSPR_INPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_OUTPUT_CHANNELS", customEnvValues)
	os.Setenv("INSPR_UNIX_SOCKET", unixSocketAddr)
	os.Setenv("INSPR_APP_CTX", "random.ctx")
	os.Setenv("INSPR_ENV", "test")
	os.Setenv("INSPR_APP_ID", "appid")
}

// deleteMockEnvVars - deletes the env values used in the tests functions
func deleteMockEnvVars() {
	os.Unsetenv("INSPR_OUTPUT_CHANNELS")
	os.Unsetenv("INSPR_INPUT_CHANNELS")
	os.Unsetenv("INSPR_UNIX_SOCKET")
	os.Unsetenv("INSPR_APP_CTX")
	os.Unsetenv("INSPR_ENV")
	os.Unsetenv("INSPR_APP_ID")
}

type mockReader struct {
	readMessage func(ctx context.Context, channel string) (models.BrokerData, error)
	commit      func(ctx context.Context, channel string) error
}

func (m mockReader) Commit(ctx context.Context, channel string) error {
	return m.commit(ctx, channel)
}
func (m mockReader) ReadMessage(ctx context.Context, channel string) (models.BrokerData, error) {
	return m.readMessage(ctx, channel)
}

type mockWriter struct {
	writeMessage func(channel string, message interface{}) error
}

func (m mockWriter) WriteMessage(channel string, message interface{}) error {
	return m.writeMessage(channel, message)
}
func TestServer_writeMessageHandler(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()
	tests := []struct {
		readerFunc func(t *testing.T) mockReader
		writerFunc func(t *testing.T) mockWriter
		channel    string
		name       string
		wantErr    bool
		message    interface{}
	}{
		{
			name:    "correct behaviour test",
			channel: "chan",
			message: struct {
				Message interface{} `json:"message"`
			}{"lofi nordeste"},
		},
		{
			name:    "invalid channel",
			channel: "invalid",
			message: struct {
				Message interface{} `json:"message"`
			}{"lofi nordeste"},
			wantErr: true,
		},
		{
			name:    "invalid data for marshalling",
			channel: "chan",
			message: "invalid message",
			wantErr: true,
		},
		{
			name:    "invalid broker response",
			channel: "chan",
			message: struct {
				Message interface{} `json:"message"`
			}{"this is an invalid message"},
			wantErr: true,
			writerFunc: func(t *testing.T) mockWriter {
				return mockWriter{
					writeMessage: func(channel string, message interface{}) error {
						return errors.New("this is an error")
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.writerFunc == nil {
				tt.writerFunc = func(t *testing.T) mockWriter {
					return mockWriter{
						writeMessage: func(channel string, message interface{}) error {
							if !reflect.DeepEqual(message, tt.message.(struct {
								Message interface{} `json:"message"`
							}).Message) {

								t.Errorf("Server_writeMessageHandler message = %v, want = %v", message, tt.message)
							}

							return nil
						},
					}
				}
			}
			s := &Server{
				Writer: tt.writerFunc(t),
			}
			server := httptest.NewServer(s.writeMessageHandler())
			defer server.Close()
			client := request.NewJSONClient(server.URL)

			rest := struct {
				Status string `json:"status"`
			}{}

			err := client.Send(context.Background(), tt.channel, http.MethodPost, tt.message, &rest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server_writeMessageHandler err = %v, wantErr = %v", err, tt.wantErr)
			}

		})
	}
}

func TestServer_readMessageRoutine(t *testing.T) {
	createMockEnvVars()
	defer deleteMockEnvVars()
	channels := map[string]chan interface{}{
		"chan":    make(chan interface{}, 2),
		"banana":  make(chan interface{}, 2),
		"testing": make(chan interface{}, 2),
	}
	readerFuncErr := func(t *testing.T) mockReader {
		return mockReader{
			readMessage: func(ctx context.Context, channel string) (models.BrokerData, error) {
				return models.BrokerData{}, errors.New("this is an error")
			},
		}
	}
	readerFunc := func(t *testing.T) mockReader {
		return mockReader{
			readMessage: func(ctx context.Context, channel string) (models.BrokerData, error) {
				var msg interface{}
				select {

				case msg = <-channels[channel]:
				case <-ctx.Done():
					return models.BrokerData{}, ctx.Err()
				}
				return models.BrokerData{
					Message: msg,
				}, nil
			},
			commit: func(ctx context.Context, channel string) error {
				return nil
			},
		}
	}
	tests := []struct {
		readerFunc func(t *testing.T) mockReader
		writerFunc func(t *testing.T) mockWriter
		channel    string
		name       string
		wantErr    bool
		message    interface{}
	}{

		{
			name:       "correct functionality",
			channel:    "chan",
			readerFunc: readerFunc,
		},
		{
			name:       "correct functionality",
			channel:    "testing",
			readerFunc: readerFunc,
		},
		{
			name:       "correct functionality",
			channel:    "banana",
			readerFunc: readerFunc,
		},
		{
			name:       "incorrect functionality",
			channel:    "banana",
			wantErr:    true,
			readerFunc: readerFuncErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &Server{
				Reader: tt.readerFunc(t),
			}
			received := false
			server := httptest.NewServer(rest.Handler(
				func(w http.ResponseWriter, r *http.Request) {

					received = true
					channel := strings.TrimPrefix(r.URL.Path, "/")
					if channel != tt.channel {
						t.Errorf("Server_readMessageRoutine %v = %v , want %v", "channel", channel, tt.channel)
					}

					decoder := json.NewDecoder(r.Body)
					ret := struct {
						Message interface{} `json:"message"`
					}{}
					err := decoder.Decode(&ret)
					if (err != nil) != tt.wantErr {
						t.Errorf("Server_readMessageRoutine err = %v, wantErr = %v", err, tt.wantErr)
					}

					if !reflect.DeepEqual(ret.Message, tt.message) {
						t.Errorf("Server_readMessageRoutine message = %v, want %v", ret.Message, tt.message)
					}
					encoder := json.NewEncoder(w)
					encoder.Encode(struct {
						Status string `json:"status"`
					}{"OK"})
				},
			))
			defer server.Close()

			s.client = request.NewJSONClient(server.URL)
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
			defer cancel()
			go s.readMessageRoutine(ctx)

			channels[tt.channel] <- tt.message
			select {

			case <-ctx.Done():
				if received == tt.wantErr {
					t.Errorf("Server_readMessageRoutine received = %v, wantErr = %v", received, tt.wantErr)
				}
			}
		})
	}
}
