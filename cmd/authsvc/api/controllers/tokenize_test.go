package controllers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/auth"
)

const bitSize = 512

func generatePrivateKey() (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func TestServer_Tokenize(t *testing.T) {
	tests := []struct {
		name string
		want int
		body auth.Payload
	}{
		{
			name: "Tokenize_valid_payload",
			want: http.StatusOK,
			body: auth.Payload{
				UID:        "u000001",
				Scope:      []string{""},
				Role:       1,
				Refresh:    "refreshtk",
				RefreshURL: "http://refresh.token",
			},
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
			handlerFunc := server.Tokenize().HTTPHandlerFunc()
			ts := httptest.NewServer(handlerFunc)
			defer ts.Close()

			client := ts.Client()
			body, err := json.Marshal(tt.body)
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
			}
		})
	}
}
