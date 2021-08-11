package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/rest/request"
	"inspr.dev/inspr/pkg/utils"
)

func TestBrokersClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *models.BrokersDI
		wantErr bool
	}{
		{
			name: "get broker controler test",
			args: args{
				ctx: context.Background(),
			},
			want: &models.BrokersDI{
				Available: utils.StringArray{"mock_broker"},
				Default:   "mock_broker",
			},
			wantErr: false,
		},
		{
			name: "failed get broker controler test",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				encoder := json.NewEncoder(w)
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					encoder.Encode(ierrors.New("").BadRequest())
					return
				}

				if r.URL.Path != "/brokers" {
					t.Errorf("path is not brokers")
				}

				if r.Method != http.MethodGet {
					t.Errorf("method is not GET")
				}

				di := models.BrokersDI{
					Available: utils.StringArray{"mock_broker"},
					Default:   "mock_broker",
				}

				encoder.Encode(di)
			}
			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			cc := &BrokersClient{
				reqClient: request.NewJSONClient(s.URL),
			}
			got, err := cc.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("BrokersClient.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BrokersClient.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrokersClient_Create(t *testing.T) {
	type args struct {
		ctx        context.Context
		brokerName string
		config     []byte
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "send_request_to_route",
			args: args{
				ctx:        context.Background(),
				brokerName: "kafka",
				config:     []byte{},
			},
			want: nil,
		},
		{
			name: "failed_to_send_request_to_route",
			args: args{
				ctx:        context.Background(),
				brokerName: "kafka",
				config:     []byte{},
			},
			want: ierrors.New("error_on_request"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				if tt.want != nil {
					w.WriteHeader(http.StatusBadRequest)
					rest.ERROR(w, ierrors.New(tt.want).BadRequest())
					return
				}

				if r.URL.Path != "/brokers/"+tt.args.brokerName {
					t.Errorf("path is not brokers/" + tt.args.brokerName)
				}

				if r.Method != http.MethodPost {
					t.Errorf("method is not " + http.MethodPost)
				}
				rest.JSON(w, http.StatusOK, nil)
			}

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()

			bc := &BrokersClient{
				reqClient: request.NewJSONClient(s.URL),
			}

			err := bc.Create(tt.args.ctx, tt.args.brokerName, tt.args.config)

			fmt.Println(err, tt.want)

			var got string
			if err != nil {
				got = err.Error()
			}

			if (tt.want != nil) && got != tt.want.Error() {
				t.Errorf("BrokersClient.Create() error = %v, want %v", err, tt.want)
			}
		})
	}
}
