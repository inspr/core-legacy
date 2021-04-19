package auth

import (
	"crypto/rsa"
	"encoding/pem"
	"os"

	"github.com/inspr/inspr/pkg/ierrors"
	"golang.org/x/crypto/ssh"
)

func GetPublicKey() (*rsa.PublicKey, error) {
	pubBytes, ok := os.LookupEnv("AUTH_PATH")
	if !ok {
		err := ierrors.NewError().Message("AUTH_PATH unavailible").Build()
		return nil, err
	}

	pubBlock, _ := pem.Decode([]byte(pubBytes))
	if pubBlock.Type != "RSA PRIVATE KEY" {
		err := ierrors.NewError().InternalServer().Message("RSA public key is of the wrong type").Build()
		return nil, err
	}

	parsed, _, _, _, err := ssh.ParseAuthorizedKey(pubBlock.Bytes)
	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)
	if err != nil {
		return nil, err
	}

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	pub := pubCrypto.(*rsa.PublicKey)
	return pub, nil
}
