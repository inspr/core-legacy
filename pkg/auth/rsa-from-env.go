package auth

import (
	"crypto/rsa"
	"encoding/pem"
	"os"

	"golang.org/x/crypto/ssh"
	"inspr.dev/inspr/pkg/ierrors"
)

// GetPublicKey resolves the rsa public key from the environment variable
func GetPublicKey() (*rsa.PublicKey, error) {
	pubBytes, ok := os.LookupEnv("JWT_PUBLIC_KEY")
	if !ok {
		err := ierrors.NewError().Message("JWT_PUBLIC_KEY unavailable").Build()
		return nil, err
	}

	pubBlock, _ := pem.Decode([]byte(pubBytes))
	if pubBlock.Type != "RSA PUBLIC KEY" {
		err := ierrors.NewError().InternalServer().Message("RSA public key is of the wrong type").Build()
		return nil, err
	}

	parsed, _, _, _, err := ssh.ParseAuthorizedKey(pubBlock.Bytes)
	if err != nil {
		return nil, err
	}

	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)

	// Then, CryptoPublicKey() is called to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, the result is converted back to a *rsa.PublicKey
	pub := pubCrypto.(*rsa.PublicKey)
	return pub, nil
}
