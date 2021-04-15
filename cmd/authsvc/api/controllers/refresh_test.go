package controllers

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/rest"
)

func TestServer_Refresh(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		payload models.Payload
		body    func(string) interface{}
		status  int
	}{
		{
			name: "Valid refresh",
			want: http.StatusOK,
			payload: models.Payload{
				UID:        "u000001",
				Scope:      []string{""},
				Role:       1,
				Refresh:    []byte("refreshtk"),
				RefreshURL: "http://refresh.token",
			},
			body: func(s string) interface{} {
				return models.ResfreshDI{
					RefreshToken: []byte("mock_token"),
					RefreshURL:   s,
				}
			},
			status: 200,
		},
		{
			name: "Invalid refresh, bad request",
			want: http.StatusBadRequest,
			body: func(s string) interface{} {
				return struct {
					RefreshURL string `json:"refreshtoken"`
					Data       int    `json:"refreshurl"`
				}{
					Data:       0,
					RefreshURL: s,
				}
			},
		},
		{
			name: "Invalid refresh, UID refresh fail",
			want: http.StatusInternalServerError,
			body: func(s string) interface{} {
				return models.ResfreshDI{
					RefreshToken: []byte("mock_token"),
					RefreshURL:   s,
				}
			},
			status: http.StatusInternalServerError,
		},
	}

	privKey, _ := generatePrivateKey()
	privDER := x509.MarshalPKCS1PrivateKey(privKey)
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privPEM := pem.EncodeToMemory(&privBlock)

	// Configuring enviroment for tests
	os.Setenv("JWT_PRIVATE_KEY", string(privPEM))

	var server Server
	server.Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc := server.Refresh().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			auxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rest.JSON(w, tt.status, tt.payload)
			})
			auxServer := httptest.NewServer(auxHandler)
			defer auxServer.Close()

			client := ts.Client()

			body, err := json.Marshal(tt.body(auxServer.URL))

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
				t.Errorf("AuthHandlers_Tokenize() = %v, want %v", res.StatusCode, tt.want)
				return
			}

			if tt.want == 200 {
				jwtdo := models.JwtDO{}
				err = json.NewDecoder(res.Body).Decode(&jwtdo)
				if err != nil {
					t.Log("error making a POST in the httptest server")
					return
				}

				payload, err := auth.Desserialize(jwtdo.Token)
				if err != nil {
					t.Errorf("AuthHandlers_Tokenize(), %v", err.Error())
					return
				}

				if !reflect.DeepEqual(payload, tt.payload) {
					t.Errorf("AuthHandlers_Tokenize() = %v, want %v", payload, tt.payload)
					return
				}
			}
		})
	}
}
