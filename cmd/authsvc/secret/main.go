package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientSet kubernetes.Interface

const bitSize = 256

// initKube initializes a k8s operator with in cluster configuration
func initKube() error {

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

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

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

func main() {
	namespace := os.Getenv("CURR_NAMESPACE")

	initKube()

	_, errPriv := clientSet.CoreV1().Secrets(namespace).Get("jwtprivatekey", v1.GetOptions{})
	_, errPub := clientSet.CoreV1().Secrets(namespace).Get("jwtpublickey", v1.GetOptions{})

	if errPriv != nil || errPub != nil {
		if errPriv == nil {
			clientSet.CoreV1().Secrets(namespace).Delete("jwtprivatekey", &v1.DeleteOptions{})
		}

		if errPub == nil {
			clientSet.CoreV1().Secrets(namespace).Delete("jwtpublickey", &v1.DeleteOptions{})
		}

		privateKey, err := generatePrivateKey()
		if err != nil {
			log.Fatal(err.Error())
		}

		publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
		if err != nil {
			log.Fatal(err.Error())
		}

		privateKeyBytes := encodePrivateKeyToPEM(privateKey)

		privSec := corev1.Secret{
			Type: corev1.SecretTypeOpaque,
			ObjectMeta: v1.ObjectMeta{
				Name: "jwtprivatekey",
			},
			Data: map[string][]byte{
				"key": privateKeyBytes,
			},
		}
		pubSec := corev1.Secret{
			Type: corev1.SecretTypeOpaque,
			ObjectMeta: v1.ObjectMeta{
				Name: "jwtpublickey",
			},
			Data: map[string][]byte{
				"key": publicKeyBytes,
			},
		}
		_, err = clientSet.CoreV1().Secrets(namespace).Create(&privSec)
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = clientSet.CoreV1().Secrets(namespace).Create(&pubSec)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("New secrets generated.")
	}
}
