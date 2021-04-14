package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientSet kubernetes.Interface

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

func main() {

	namespace := os.Getenv("K8S_NAMESPACE")
	pvtKeyName := os.Getenv("PVT_KEY_NAME")
	if pvtKeyName == "" {
		panic("[ENV VAR] PVT_KEY_NAME not found")
	}

	initKube()

	_, err := clientSet.CoreV1().Secrets(namespace).Get("redisprivatekey", v1.GetOptions{})

	if err != nil {
		bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
		if _, err := rand.Read(bytes); err != nil {
			panic(err.Error())
		}

		key := hex.EncodeToString(bytes)
		privateKeyBytes := []byte(key)

		privSec := corev1.Secret{
			Type: corev1.SecretTypeOpaque,
			ObjectMeta: v1.ObjectMeta{
				Name: pvtKeyName,
			},
			Data: map[string][]byte{
				"key": privateKeyBytes,
			},
		}

		_, err = clientSet.CoreV1().Secrets(namespace).Create(&privSec)
		if err != nil {
			log.Fatal(err.Error())
		}

	}
}
