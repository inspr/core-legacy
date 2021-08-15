package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestGetPublicKey(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 512)
	err := privateKey.Validate()
	if err != nil {
		t.Errorf("Could't generate valid key for testing pourpuses")
	}
	publicRsaKey, _ := ssh.NewPublicKey(&privateKey.PublicKey)
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	tests := []struct {
		name     string
		env      bool
		pemType  string
		keyBytes []byte
		wantErr  bool
	}{
		{
			name:    "Invalid environment variable",
			env:     false,
			wantErr: true,
		},
		{
			name:     "Invalid environment pem type",
			env:      true,
			pemType:  "WRONG",
			keyBytes: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid environment pem key",
			env:      true,
			pemType:  "RSA PUBLIC KEY",
			keyBytes: []byte{45, 78, 56},
			wantErr:  true,
		},
		{
			name:     "Valid",
			env:      true,
			pemType:  "RSA PUBLIC KEY",
			keyBytes: publicKeyBytes,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env {
				pubBlock := pem.Block{
					Type:    tt.pemType,
					Headers: nil,
					Bytes:   tt.keyBytes,
				}
				publicPEM := pem.EncodeToMemory(&pubBlock)
				os.Setenv("JWT_PUBLIC_KEY", string(publicPEM))
			}
			got, err := GetPublicKey()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"GetPublicKey() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr && !privateKey.PublicKey.Equal(got) {
				t.Errorf(
					"GetPublicKey() = %v, want %v",
					got,
					privateKey.PublicKey,
				)
			}
		})
	}
}
