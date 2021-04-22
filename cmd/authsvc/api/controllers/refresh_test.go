package controllers

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/auth/models"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

var server Server

func TestServer_Refresh(t *testing.T) {
	tests := []struct {
		name      string
		want      int
		payload   models.Payload
		importURL bool
		headToken bool
		token     func(models.Payload) []byte
		status    int
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
			token:     getToken,
			headToken: true,
			importURL: true,
			status:    200,
		},
		{
			name: "Invalid refresh, header missing 'Authentication'",
			want: http.StatusUnauthorized,
			token: func(models.Payload) []byte {
				return nil
			},
			headToken: false,
		},
		{
			name: "Invalid refresh, invalid token signature",
			want: http.StatusForbidden,
			token: func(models.Payload) []byte {
				return []byte("cicada")
			},
			headToken: true,
		},
		{
			name: "Invalid refresh, invalid token payload",
			want: http.StatusForbidden,
			token: func(models.Payload) []byte {
				token := jwt.New()
				payload := struct {
					UID int
				}{
					UID: 3301,
				}
				token.Set("payload", payload)
				signed, _ := jwt.Sign(token, jwa.RS256, server.privKey)
				return signed
			},
			headToken: true,
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

	server.Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test server for refresh handler
			ts := httptest.NewServer(server.Mux)
			defer ts.Close()

			// Test server for mocking UID server
			mockUIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rest.JSON(w, tt.status, tt.payload)
			})
			mockUIDServer := httptest.NewServer(mockUIDHandler)
			defer mockUIDServer.Close()

			// Token for request header
			if tt.importURL {
				tt.payload.RefreshURL = mockUIDServer.URL
			}
			signed := tt.token(tt.payload)

			// Request assembly
			req, _ := http.NewRequest(http.MethodGet, ts.URL+"/refresh", nil)
			head := http.Header{}
			if tt.headToken {
				head.Add("Authorization", fmt.Sprintf("Bearer %v", string(signed)))
			}
			req.Header = head

			testClient := ts.Client()

			res, err := testClient.Do(req)
			if err != nil {
				t.Errorf("error making a GET in the httptest server, error: %s", err.Error())
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

				if !reflect.DeepEqual(*payload, tt.payload) {
					t.Errorf("AuthHandlers_Tokenize() = %v, want %v", payload, tt.payload)
					return
				}
			}
		})
	}
}

func getToken(payload models.Payload) []byte {
	token := jwt.New()
	token.Set("payload", payload)
	signed, _ := jwt.Sign(token, jwa.RS256, server.privKey)
	return signed
}
