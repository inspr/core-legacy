package auth

import (
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/inspr/inspr/pkg/ierrors"
	"golang.org/x/crypto/ssh"
)

// GetPublicKey resolves the rsa public key from the environment variable
func GetPublicKey() (*rsa.PublicKey, error) {
	pubBytes, ok := os.LookupEnv("JWT_PUBLIC_KEY")
	if !ok {
		err := ierrors.NewError().Message("JWT_PUBLIC_KEY unavailable").Build()
		return nil, err
	}
	fmt.Printf("%s\n", pubBytes)

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

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	pub := pubCrypto.(*rsa.PublicKey)
	return pub, nil
}
