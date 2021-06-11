package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/inspr/inspr/pkg/utils"
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
				Installed: utils.StringArray{"mock_broker"},
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
					encoder.Encode(ierrors.NewError().BadRequest().Build())
					return
				}

				if r.URL.Path != "/brokers" {
					t.Errorf("path is not channels")
				}

				if r.Method != http.MethodGet {
					t.Errorf("method is not GET")
				}

				di := models.BrokersDI{
					Installed: utils.StringArray{"mock_broker"},
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
