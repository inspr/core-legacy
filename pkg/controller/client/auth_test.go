package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/api/models"
	auth "github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/rest/request"
)

func TestAuthClient_GenerateToken(t *testing.T) {
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
			name: "generate token test",
			args: args{
				ctx: context.Background(),
				payload: auth.Payload{
					UID:        "test123",
					Permission: nil,
					Scope:      []string{"app1", "app2"},
					Refresh:    []byte("refreshtoken1234"),
					RefreshURL: "http://URLToUIDProvider.valid",
				},
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "generate token with error test",
			args: args{
				ctx:     context.Background(),
				payload: auth.Payload{},
			},
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

				if r.URL.Path != "/auth" {
					t.Errorf("path is not auth")
				}

				if r.Method != "POST" {
					t.Errorf("method is not POST")
				}

				var payload auth.Payload
				var authDI models.AuthDI

				decoder := request.JSONDecoderGenerator(r.Body)
				err := decoder.Decode(&payload)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(payload, tt.args.payload) {
					t.Errorf("wrong token. want = %v, got = %v", payload, tt.args.payload)
				}

				encoder.Encode(authDI)
				tt.want = string(authDI.Token)
			}

			s := httptest.NewServer(http.HandlerFunc(handler))
			defer s.Close()
			ac := &AuthClient{
				c: request.NewJSONClient(s.URL),
			}
			got, err := ac.GenerateToken(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthorizationClient.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthorizationClient.GenerateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
