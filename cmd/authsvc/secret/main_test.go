package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

func Test_KeyGen(t *testing.T) {
	logger, _ = zap.NewDevelopment(zap.Fields(zap.String("section", "Auth-provider")))

	privateKey, err := generatePrivateKey()
	if err != nil {
		logger.Fatal(err.Error())
	}

	privateKeyBytes, publicKeyBytes, err := encodeKeysToPEM(privateKey)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if ok := verifyKeyPair(string(privateKeyBytes), string(publicKeyBytes)); !ok {
		t.Errorf("alalala")
	}

}

func verifyKeyPair(private, public string) bool {
	block, _ := pem.Decode([]byte(private))
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	pubBlock, _ := pem.Decode([]byte(public))
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(pubBlock.Bytes)
	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	pub := pubCrypto.(*rsa.PublicKey)
	return key.PublicKey.Equal(pub)
}
